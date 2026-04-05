package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
)

// === Models ===

type BookingSchedule struct {
	ID string `json:"id"`

	OriginAreaID string `json:"origin_area_id"`
	OriginName   string `json:"origin_name"`
	DestAreaID   string `json:"dest_area_id"`
	DestName     string `json:"dest_name"`

	TravelDate string `json:"travel_date"`
	TimeFrom   string `json:"time_from"`
	TimeTo     string `json:"time_to"`

	SeatType  string `json:"seat_type"`
	SeatCount int    `json:"seat_count"`
	AutoBook  bool   `json:"auto_book"`

	PassengerName  string `json:"passenger_name"`
	PassengerPhone string `json:"passenger_phone"`
	PassengerEmail string `json:"passenger_email"`

	Status string `json:"status"`

	BookingID     string     `json:"booking_id"`
	BookingCode   string     `json:"booking_code"`
	TicketPrice   int        `json:"ticket_price"`
	SeatName      string     `json:"seat_name"`
	DepartureTime *time.Time `json:"departure_time"`

	RetryCount int    `json:"retry_count"`
	LastError  string `json:"last_error"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	RouteName string `json:"route_name,omitempty"`
}

type Stats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Searching int `json:"searching"`
	Success   int `json:"success"`
	Failed    int `json:"failed"`
	Cancelled int `json:"cancelled"`
}

// === DB ===

type DB struct {
	Pool *pgxpool.Pool
}

func New(cfg config.DatabaseConfig) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

// === Schedules ===

const scheduleColumns = `id,
	origin_area_id, origin_name,
	dest_area_id, dest_name,
	travel_date, time_from, time_to,
	seat_type, seat_count, auto_book,
	passenger_name, passenger_phone, passenger_email,
	status, booking_id, booking_code, ticket_price,
	seat_name, departure_time,
	retry_count, last_error,
	created_at, updated_at, route_name`

func scanSchedule(scan func(dest ...any) error) (*BookingSchedule, error) {
	var s BookingSchedule
	var travelDate time.Time
	err := scan(
		&s.ID,
		&s.OriginAreaID, &s.OriginName,
		&s.DestAreaID, &s.DestName,
		&travelDate, &s.TimeFrom, &s.TimeTo,
		&s.SeatType, &s.SeatCount, &s.AutoBook,
		&s.PassengerName, &s.PassengerPhone, &s.PassengerEmail,
		&s.Status, &s.BookingID, &s.BookingCode, &s.TicketPrice,
		&s.SeatName, &s.DepartureTime,
		&s.RetryCount, &s.LastError,
		&s.CreatedAt, &s.UpdatedAt,
		&s.RouteName,
	)
	if err != nil {
		return nil, err
	}
	s.TravelDate = travelDate.Format("2006-01-02")
	return &s, nil
}

func (db *DB) CreateSchedule(ctx context.Context, s *BookingSchedule) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO booking_schedules (
			origin_area_id, origin_name,
			dest_area_id, dest_name,
			travel_date, time_from, time_to,
			seat_type, seat_count, auto_book,
			passenger_name, passenger_phone, passenger_email,
			status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,'pending')
		RETURNING id, created_at, updated_at`,
		s.OriginAreaID, s.OriginName,
		s.DestAreaID, s.DestName,
		s.TravelDate, s.TimeFrom, s.TimeTo,
		s.SeatType, s.SeatCount, s.AutoBook,
		s.PassengerName, s.PassengerPhone, s.PassengerEmail,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (db *DB) ListSchedules(ctx context.Context, statusFilter string) ([]BookingSchedule, error) {
	return db.listSchedules(ctx, "", statusFilter)
}

func (db *DB) ListSchedulesByEmail(ctx context.Context, emailAddr, statusFilter string) ([]BookingSchedule, error) {
	return db.listSchedules(ctx, emailAddr, statusFilter)
}

func (db *DB) listSchedules(ctx context.Context, emailAddr, statusFilter string) ([]BookingSchedule, error) {
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules`
	args := []any{}
	conditions := []string{}
	argIdx := 1

	if emailAddr != "" {
		conditions = append(conditions, fmt.Sprintf("passenger_email = $%d", argIdx))
		args = append(args, emailAddr)
		argIdx++
	}

	if statusFilter != "" && statusFilter != "all" {
		switch statusFilter {
		case "active":
			conditions = append(conditions, "status IN ('pending','searching','found','booking')")
		case "success":
			conditions = append(conditions, "status = 'success'")
		case "failed":
			conditions = append(conditions, "status IN ('failed','cancelled','expired')")
		default:
			conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
			args = append(args, statusFilter)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for _, c := range conditions[1:] {
			query += " AND " + c
		}
	}

	query += ` ORDER BY created_at DESC`

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BookingSchedule
	for rows.Next() {
		s, err := scanSchedule(rows.Scan)
		if err != nil {
			return nil, err
		}
		results = append(results, *s)
	}
	return results, nil
}

func (db *DB) GetSchedule(ctx context.Context, id string) (*BookingSchedule, error) {
	row := db.Pool.QueryRow(ctx,
		`SELECT `+scheduleColumns+` FROM booking_schedules WHERE id = $1`, id)
	return scanSchedule(row.Scan)
}

func (db *DB) DeleteSchedule(ctx context.Context, id string) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM booking_schedules WHERE id = $1`, id)
	return err
}

func (db *DB) CancelSchedule(ctx context.Context, id string) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE booking_schedules SET status='cancelled', updated_at=NOW() WHERE id=$1 AND status NOT IN ('success','cancelled')`, id)
	return err
}

func (db *DB) GetStats(ctx context.Context) (*Stats, error) {
	return db.getStats(ctx, "")
}

func (db *DB) GetStatsByEmail(ctx context.Context, emailAddr string) (*Stats, error) {
	return db.getStats(ctx, emailAddr)
}

func (db *DB) getStats(ctx context.Context, emailAddr string) (*Stats, error) {
	var s Stats
	where := ""
	args := []any{}
	if emailAddr != "" {
		where = " WHERE passenger_email = $1"
		args = append(args, emailAddr)
	}
	err := db.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status IN ('searching','found','booking')),
			COUNT(*) FILTER (WHERE status = 'success'),
			COUNT(*) FILTER (WHERE status IN ('failed','expired')),
			COUNT(*) FILTER (WHERE status = 'cancelled')
		FROM booking_schedules`+where,
		args...,
	).Scan(&s.Total, &s.Pending, &s.Searching, &s.Success, &s.Failed, &s.Cancelled)
	return &s, err
}

func (db *DB) GetPendingSchedules(ctx context.Context, maxRetries int) ([]BookingSchedule, error) {
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules
		WHERE status IN ('pending', 'searching')
		AND travel_date >= CURRENT_DATE
		ORDER BY travel_date ASC, created_at ASC`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BookingSchedule
	for rows.Next() {
		s, err := scanSchedule(rows.Scan)
		if err != nil {
			return nil, err
		}
		results = append(results, *s)
	}
	return results, nil
}

func (db *DB) UpdateScheduleStatus(ctx context.Context, id, status, lastError string) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE booking_schedules SET status=$1, last_error=$2, retry_count=retry_count+1, updated_at=NOW() WHERE id=$3`,
		status, lastError, id)
	return err
}

func (db *DB) UpdateScheduleSuccess(ctx context.Context, id, bookingID, bookingCode, routeName, seatName string, price int, departureTime *time.Time) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE booking_schedules SET
			booking_id=$1, booking_code=$2, route_name=$3,
			seat_name=$4, ticket_price=$5, departure_time=$6,
			updated_at=NOW()
		WHERE id=$7`,
		bookingID, bookingCode, routeName, seatName, price, departureTime, id)
	return err
}

func (db *DB) GetRecentSchedules(ctx context.Context, limit int) ([]BookingSchedule, error) {
	return db.getRecentSchedules(ctx, "", limit)
}

func (db *DB) GetRecentSchedulesByEmail(ctx context.Context, emailAddr string, limit int) ([]BookingSchedule, error) {
	return db.getRecentSchedules(ctx, emailAddr, limit)
}

func (db *DB) getRecentSchedules(ctx context.Context, emailAddr string, limit int) ([]BookingSchedule, error) {
	where := ""
	args := []any{}
	if emailAddr != "" {
		where = " WHERE passenger_email = $1"
		args = append(args, emailAddr)
		args = append(args, limit)
	} else {
		args = append(args, limit)
	}

	limitParam := fmt.Sprintf("$%d", len(args))
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules` + where + ` ORDER BY updated_at DESC LIMIT ` + limitParam

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []BookingSchedule
	for rows.Next() {
		s, err := scanSchedule(rows.Scan)
		if err != nil {
			return nil, err
		}
		results = append(results, *s)
	}
	return results, nil
}
