package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/database"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/email"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/futa"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/webhook"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()
	if err := db.RunMigrations(context.Background(), database.ResolveMigrationsDir(cfgPath)); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	futaClient := futa.NewClient(cfg.Futa)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down worker...")
		cancel()
	}()

	log.Printf("Worker started: concurrency=%d poll_interval=%s retry_delay=%s",
		cfg.Worker.Concurrency, cfg.Worker.PollInterval, cfg.Worker.RetryDelay)

	var wg sync.WaitGroup
	for i := 0; i < cfg.Worker.Concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorkerLoop(ctx, id, db, futaClient, cfg)
		}(i)
	}
	wg.Wait()
	log.Println("Worker stopped")
}

// runWorkerLoop is the main loop for a single worker goroutine. It continuously
// tries to claim and process one job at a time. When no job is available it
// backs off for PollInterval before trying again.
func runWorkerLoop(ctx context.Context, id int, db *database.DB, client *futa.Client, cfg *config.Config) {
	wh := webhook.NewSender(cfg.Webhook)
	em := email.NewSender(cfg.Email)

	for {
		if ctx.Err() != nil {
			return
		}

		s, err := db.ClaimNextSchedule(ctx, cfg.Worker.RetryDelay)
		if err != nil {
			log.Printf("[worker-%d] ERROR claiming schedule: %v", id, err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(cfg.Worker.PollInterval):
			}
			continue
		}

		if s == nil {
			// No job available — wait before polling again.
			select {
			case <-ctx.Done():
				return
			case <-time.After(cfg.Worker.PollInterval):
			}
			continue
		}

		processOne(ctx, db, client, wh, em, *s)
	}
}

func processOne(ctx context.Context, db *database.DB, client *futa.Client, wh *webhook.Sender, em *email.Sender, s database.BookingSchedule) {
	log.Printf("[%s] Processing: %s -> %s on %s (%s-%s)",
		s.ID[:8], s.OriginName, s.DestName, s.TravelDate, s.TimeFrom, s.TimeTo)

	originAreaID := s.OriginAreaID
	destAreaID := s.DestAreaID

	// Search routes
	fromDate := s.TravelDate + "T00:00:00.000+07:00"
	routes, err := client.SearchRoutes(ctx, originAreaID, destAreaID, fromDate)
	if err != nil {
		db.UpdateScheduleStatus(ctx, s.ID, "searching", fmt.Sprintf("search routes: %v", err))
		return
	}
	if len(routes) == 0 {
		db.UpdateScheduleStatus(ctx, s.ID, "searching", "no routes found")
		return
	}

	routeIDs := make([]string, len(routes))
	for i, r := range routes {
		routeIDs[i] = r.RouteID
	}

	// Search trips
	toDate := s.TravelDate + "T23:59:59.000+07:00"
	trips, err := client.SearchTripsByRoute(ctx, routeIDs, fromDate, toDate)
	if err != nil {
		db.UpdateScheduleStatus(ctx, s.ID, "searching", fmt.Sprintf("search trips: %v", err))
		return
	}
	if len(trips) == 0 {
		db.UpdateScheduleStatus(ctx, s.ID, "searching", "no trips available")
		return
	}

	log.Printf("[%s] Found %d trips", s.ID[:8], len(trips))

	// Prefer later departures first when multiple trips match.
	sort.SliceStable(trips, func(i, j int) bool {
		a := trips[i].RawDepartureTime
		b := trips[j].RawDepartureTime
		if a != b {
			return a > b
		}
		return trips[i].DepartureTime > trips[j].DepartureTime
	})

	// Filter and find suitable trip
	for _, trip := range trips {
		if trip.EmptySeatQuantity < s.SeatCount {
			continue
		}

		// Filter by time range
		if s.TimeFrom != "" && s.TimeTo != "" && s.TimeFrom != "00:00" {
			depTime := trip.RawDepartureTime
			if depTime < s.TimeFrom || depTime > s.TimeTo {
				continue
			}
		}

		// Filter by seat type
		if s.SeatType != "any" && s.SeatType != "" {
			if s.SeatType == "giuong_nam" && trip.SeatTypeCode != "glm" {
				continue
			}
			if s.SeatType == "ghe_ngoi" && trip.SeatTypeCode == "glm" {
				continue
			}
		}

		// Filter by max price
		if s.MaxPrice > 0 && trip.Price > s.MaxPrice {
			continue
		}

		// If not auto_book, just mark as found
		if !s.AutoBook {
			db.UpdateScheduleStatus(ctx, s.ID, "found", fmt.Sprintf("Found trip %s at %s, %d empty seats", trip.TripID, trip.RawDepartureTime, trip.EmptySeatQuantity))
			return
		}

		// Get seat diagram
		seats, err := client.GetSeatDiagram(ctx, trip.TripID)
		if err != nil {
			log.Printf("[%s] Error getting seats: %v", s.ID[:8], err)
			continue
		}

		availableSeats := pickPreferredSeats(
			seats,
			s.SeatCount,
			normalizeSeatFloor(s.SeatFloor),
			normalizeSeatWindow(s.SeatWindow),
			s.PriorityTopRows,
		)

		if len(availableSeats) < s.SeatCount {
			continue
		}

		// Get departments for pickup/dropoff
		depts, err := client.GetDepartmentsInWay(ctx, trip.WayID, trip.RouteID)
		if err != nil || len(depts) < 2 {
			continue
		}

		pickup := depts[0]
		dropoff := depts[len(depts)-1]

		// Prefer origin/dest hubs
		for _, d := range depts {
			if d.PointKind == 0 {
				pickup = d
				break
			}
		}
		for i := len(depts) - 1; i >= 0; i-- {
			if depts[i].PointKind == 1 {
				dropoff = depts[i]
				break
			}
		}

		seatRefs := make([]futa.SeatRef, len(availableSeats))
		for i, seat := range availableSeats {
			seatRefs[i] = futa.SeatRef{SeatID: seat.SeatID}
		}

		booking, err := client.BookReservation(ctx,
			futa.PassengerInfo{
				CustName:    s.PassengerName,
				LoginMobile: s.PassengerPhone,
				CustEmail:   s.PassengerEmail,
				CustSn:      "",
				CustMobile:  s.PassengerPhone,
			},
			futa.TicketInfo{
				Seats:  seatRefs,
				TripID: trip.TripID,
				Pickup: futa.LocationRef{
					OfficeID:         pickup.DepartmentID,
					Name:             pickup.DepartmentName,
					Address:          pickup.DepartmentAddress,
					TimeAtDepartment: pickup.TimeAtDepartment,
					Lat:              pickup.Latitude,
					Lng:              pickup.Longitude,
					Type:             3,
				},
				Dropoff: futa.LocationRef{
					OfficeID:         dropoff.DepartmentID,
					Name:             dropoff.DepartmentName,
					Address:          dropoff.DepartmentAddress,
					TimeAtDepartment: dropoff.TimeAtDepartment,
					Lat:              dropoff.Latitude,
					Lng:              dropoff.Longitude,
					Type:             3,
				},
			},
		)
		if err != nil {
			log.Printf("[%s] Booking failed: %v", s.ID[:8], err)
			db.UpdateScheduleStatus(ctx, s.ID, "searching", fmt.Sprintf("booking failed: %v", err))
			continue
		}

		// Success!
		seatNames := ""
		for i, seat := range availableSeats {
			if i > 0 {
				seatNames += ", "
			}
			seatNames += seat.Name
		}

		depTime, _ := time.Parse(time.RFC3339, trip.DepartureTime)
		db.UpdateScheduleSuccess(ctx, s.ID,
			booking.ID, booking.Code,
			trip.Route.Name,
			seatNames,
			booking.TotalPrice, &depTime)

		log.Printf("[%s] SUCCESS! Code: %s, Price: %d, Seats: %s",
			s.ID[:8], booking.Code, booking.TotalPrice, seatNames)

		// Send payment email
		if s.PassengerEmail != "" {
			if err := em.SendPaymentLink(email.PaymentInfo{
				BookingID:   booking.ID,
				BookingCode: booking.Code,
				PhoneNumber: s.PassengerPhone,
				ToEmail:     s.PassengerEmail,
				ToName:      s.PassengerName,
				OriginName:  s.OriginName,
				DestName:    s.DestName,
				TravelDate:  s.TravelDate,
				SeatName:    seatNames,
				TicketPrice: booking.TotalPrice,
				RouteName:   trip.Route.Name,
			}); err != nil {
				log.Printf("[%s] Payment email failed: %v", s.ID[:8], err)
			}
		}

		// Send webhook
		updated, _ := db.GetSchedule(ctx, s.ID)
		if updated != nil {
			if err := wh.Send(ctx, *updated); err != nil {
				log.Printf("[%s] Webhook failed: %v", s.ID[:8], err)
			}
		}
		return
	}

	db.UpdateScheduleStatus(ctx, s.ID, "searching", "no suitable trip/seats found this round")
}

func normalizeSeatFloor(v string) string {
	switch v {
	case "up", "down", "any":
		return v
	default:
		return "any"
	}
}

func normalizeSeatWindow(v string) string {
	switch v {
	case "window", "non_window", "any":
		return v
	default:
		return "any"
	}
}

type floorMeta struct {
	minCol      int
	maxCol      int
	priorityRow map[int]bool
}

func pickPreferredSeats(seats []futa.SeatDiagramData, seatCount int, floorPref, windowPref string, priorityTopRows int) []futa.SeatDiagramData {
	if seatCount <= 0 {
		return nil
	}
	if priorityTopRows < 0 {
		priorityTopRows = 0
	}

	filtered := make([]futa.SeatDiagramData, 0, len(seats))
	for _, seat := range seats {
		if len(seat.Status) != 0 {
			continue
		}
		if floorPref != "any" && seat.Floor != floorPref {
			continue
		}
		filtered = append(filtered, seat)
	}
	if len(filtered) < seatCount {
		return nil
	}

	floorRows := map[string][]int{}
	floorRowsSet := map[string]map[int]bool{}
	meta := map[string]*floorMeta{}

	for _, seat := range filtered {
		m, ok := meta[seat.Floor]
		if !ok {
			m = &floorMeta{minCol: seat.ColumnNo, maxCol: seat.ColumnNo, priorityRow: map[int]bool{}}
			meta[seat.Floor] = m
		}
		if seat.ColumnNo < m.minCol {
			m.minCol = seat.ColumnNo
		}
		if seat.ColumnNo > m.maxCol {
			m.maxCol = seat.ColumnNo
		}

		if _, ok := floorRowsSet[seat.Floor]; !ok {
			floorRowsSet[seat.Floor] = map[int]bool{}
		}
		if !floorRowsSet[seat.Floor][seat.RowNo] {
			floorRowsSet[seat.Floor][seat.RowNo] = true
			floorRows[seat.Floor] = append(floorRows[seat.Floor], seat.RowNo)
		}
	}

	for floor, rows := range floorRows {
		sort.Ints(rows)
		if priorityTopRows <= 0 {
			continue
		}
		limit := priorityTopRows
		if limit > len(rows) {
			limit = len(rows)
		}
		m := meta[floor]
		for i := 0; i < limit; i++ {
			m.priorityRow[rows[i]] = true
		}
	}

	candidates := filtered
	if windowPref != "any" {
		next := make([]futa.SeatDiagramData, 0, len(filtered))
		for _, seat := range filtered {
			m := meta[seat.Floor]
			isWindow := seat.ColumnNo == m.minCol || seat.ColumnNo == m.maxCol
			if windowPref == "window" && isWindow {
				next = append(next, seat)
			}
			if windowPref == "non_window" && !isWindow {
				next = append(next, seat)
			}
		}
		candidates = next
	}

	if len(candidates) < seatCount {
		return nil
	}

	floorOrder := func(f string) int {
		switch f {
		case "down":
			return 0
		case "up":
			return 1
		default:
			return 2
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		a := candidates[i]
		b := candidates[j]

		aPriority := priorityTopRows > 0 && meta[a.Floor] != nil && meta[a.Floor].priorityRow[a.RowNo]
		bPriority := priorityTopRows > 0 && meta[b.Floor] != nil && meta[b.Floor].priorityRow[b.RowNo]
		if aPriority != bPriority {
			return aPriority
		}

		if floorPref == "any" && a.Floor != b.Floor {
			return floorOrder(a.Floor) < floorOrder(b.Floor)
		}
		if a.RowNo != b.RowNo {
			return a.RowNo < b.RowNo
		}
		if a.ColumnNo != b.ColumnNo {
			return a.ColumnNo < b.ColumnNo
		}
		return a.Name < b.Name
	})

	return candidates[:seatCount]
}
