package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/database"
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

	log.Printf("Worker started, polling every %s", cfg.Worker.PollInterval)

	ticker := time.NewTicker(cfg.Worker.PollInterval)
	defer ticker.Stop()

	processSchedules(ctx, db, futaClient, cfg)
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped")
			return
		case <-ticker.C:
			processSchedules(ctx, db, futaClient, cfg)
		}
	}
}

func processSchedules(ctx context.Context, db *database.DB, client *futa.Client, cfg *config.Config) {
	// Load settings for refresh token & webhook config
	settings, err := db.GetSettings(ctx)
	if err != nil {
		log.Printf("ERROR load settings: %v", err)
	} else if settings.RefreshToken != "" {
		client.SetRefreshToken(settings.RefreshToken)
	}

	// Create webhook sender from DB settings (override config)
	whCfg := cfg.Webhook
	if settings != nil && settings.WebhookURL != "" {
		whCfg.URL = settings.WebhookURL
		whCfg.Secret = settings.WebhookSecret
	}
	wh := webhook.NewSender(whCfg)

	schedules, err := db.GetPendingSchedules(ctx, cfg.Worker.MaxRetries)
	if err != nil {
		log.Printf("ERROR get pending schedules: %v", err)
		return
	}

	if len(schedules) == 0 {
		return
	}

	log.Printf("Processing %d pending schedule(s)", len(schedules))

	for _, s := range schedules {
		if ctx.Err() != nil {
			return
		}
		processOne(ctx, db, client, wh, s)
	}
}

func processOne(ctx context.Context, db *database.DB, client *futa.Client, wh *webhook.Sender, s database.BookingSchedule) {
	log.Printf("[%s] Processing: %s -> %s on %s (%s-%s)",
		s.ID[:8], s.OriginKeyword, s.DestKeyword, s.TravelDate, s.TimeFrom, s.TimeTo)

	db.UpdateScheduleStatus(ctx, s.ID, "searching", "")

	// Step 1: Resolve origin area ID
	originAreaID := s.OriginAreaID
	if originAreaID == "" {
		groups, areas, err := client.SearchPickupPoints(ctx, s.OriginKeyword)
		if err != nil {
			db.UpdateScheduleStatus(ctx, s.ID, "searching", fmt.Sprintf("search origin: %v", err))
			return
		}
		if len(areas) > 0 {
			originAreaID = areas[0].ID
		} else if len(groups) > 0 && len(groups[0].Group) > 0 {
			originAreaID = groups[0].Group[0].ProvinceID
		}
		if originAreaID == "" {
			db.UpdateScheduleStatus(ctx, s.ID, "searching", "no origin area found")
			return
		}
	}

	// Step 2: Resolve dest area ID
	destAreaID := s.DestAreaID
	if destAreaID == "" {
		groups, areas, err := client.SearchPickupPoints(ctx, s.DestKeyword)
		if err != nil {
			db.UpdateScheduleStatus(ctx, s.ID, "searching", fmt.Sprintf("search dest: %v", err))
			return
		}
		if len(areas) > 0 {
			destAreaID = areas[0].ID
		} else if len(groups) > 0 && len(groups[0].Group) > 0 {
			destAreaID = groups[0].Group[0].ProvinceID
		}
		if destAreaID == "" {
			db.UpdateScheduleStatus(ctx, s.ID, "searching", "no dest area found")
			return
		}
	}

	// Step 3: Search routes
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

	// Step 4: Search trips
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

	// Step 5: Filter and find suitable trip
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

		// If not auto_book, just mark as found
		if !s.AutoBook {
			db.UpdateScheduleStatus(ctx, s.ID, "found", fmt.Sprintf("Found trip %s at %s, %d empty seats", trip.TripID, trip.RawDepartureTime, trip.EmptySeatQuantity))
			return
		}

		// Step 6: Get seat diagram
		seats, err := client.GetSeatDiagram(ctx, trip.TripID)
		if err != nil {
			log.Printf("[%s] Error getting seats: %v", s.ID[:8], err)
			continue
		}

		var availableSeats []futa.SeatDiagramData
		for _, seat := range seats {
			if len(seat.Status) == 0 {
				availableSeats = append(availableSeats, seat)
			}
			if len(availableSeats) >= s.SeatCount {
				break
			}
		}

		if len(availableSeats) < s.SeatCount {
			continue
		}

		// Step 7: Get departments for pickup/dropoff
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

		// Step 8: Book
		db.UpdateScheduleStatus(ctx, s.ID, "booking", "")

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
			seatNames, trip.Route.Name,
			booking.TotalPrice, &depTime)

		log.Printf("[%s] SUCCESS! Code: %s, Price: %d, Seats: %s",
			s.ID[:8], booking.Code, booking.TotalPrice, seatNames)

		// Send webhook
		updated, _ := db.GetSchedule(ctx, s.ID)
		if updated != nil {
			if err := wh.Send(ctx, *updated); err != nil {
				log.Printf("[%s] Webhook failed: %v", s.ID[:8], err)
			} else {
				db.MarkWebhookSent(ctx, s.ID)
			}
		}
		return
	}

	db.UpdateScheduleStatus(ctx, s.ID, "searching", "no suitable trip/seats found this round")
}
