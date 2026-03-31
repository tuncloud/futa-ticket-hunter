package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
)

// === Models ===

type AppSettings struct {
	ID            int       `json:"id"`
	RefreshToken  string    `json:"refresh_token"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	WebhookURL    string    `json:"webhook_url"`
	WebhookSecret string    `json:"webhook_secret"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookingSchedule struct {
	ID string `json:"id"`

	OriginKeyword string `json:"origin_keyword"`
	OriginAreaID  string `json:"origin_area_id"`
	OriginName    string `json:"origin_name"`
	DestKeyword   string `json:"dest_keyword"`
	DestAreaID    string `json:"dest_area_id"`
	DestName      string `json:"dest_name"`

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

	BookingID     string          `json:"booking_id"`
	BookingCode   string          `json:"booking_code"`
	TicketPrice   int             `json:"ticket_price"`
	SeatName      string          `json:"seat_name"`
	RouteName     string          `json:"route_name"`
	DepartureTime *time.Time      `json:"departure_time"`
	TripInfo      json.RawMessage `json:"trip_info"`

	RetryCount  int    `json:"retry_count"`
	LastError   string `json:"last_error"`
	WebhookSent bool   `json:"webhook_sent"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// === Settings ===

func (db *DB) GetSettings(ctx context.Context) (*AppSettings, error) {
	var s AppSettings
	err := db.Pool.QueryRow(ctx,
		`SELECT id, refresh_token, full_name, phone, email, webhook_url, webhook_secret, updated_at
		 FROM app_settings WHERE id = 1`).Scan(
		&s.ID, &s.RefreshToken, &s.FullName, &s.Phone, &s.Email,
		&s.WebhookURL, &s.WebhookSecret, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (db *DB) UpdateSettings(ctx context.Context, s *AppSettings) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE app_settings SET
			refresh_token=$1, full_name=$2, phone=$3, email=$4,
			webhook_url=$5, webhook_secret=$6, updated_at=NOW()
		WHERE id=1`,
		s.RefreshToken, s.FullName, s.Phone, s.Email,
		s.WebhookURL, s.WebhookSecret)
	return err
}

// === Schedules ===

const scheduleColumns = `id,
	origin_keyword, origin_area_id, origin_name,
	dest_keyword, dest_area_id, dest_name,
	travel_date, time_from, time_to,
	seat_type, seat_count, auto_book,
	passenger_name, passenger_phone, passenger_email,
	status, booking_id, booking_code, ticket_price,
	seat_name, route_name, departure_time, trip_info,
	retry_count, last_error, webhook_sent,
	created_at, updated_at`

func scanSchedule(scan func(dest ...any) error) (*BookingSchedule, error) {
	var s BookingSchedule
	var travelDate time.Time
	err := scan(
		&s.ID,
		&s.OriginKeyword, &s.OriginAreaID, &s.OriginName,
		&s.DestKeyword, &s.DestAreaID, &s.DestName,
		&travelDate, &s.TimeFrom, &s.TimeTo,
		&s.SeatType, &s.SeatCount, &s.AutoBook,
		&s.PassengerName, &s.PassengerPhone, &s.PassengerEmail,
		&s.Status, &s.BookingID, &s.BookingCode, &s.TicketPrice,
		&s.SeatName, &s.RouteName, &s.DepartureTime, &s.TripInfo,
		&s.RetryCount, &s.LastError, &s.WebhookSent,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	s.TravelDate = travelDate.Format("2006-01-02")
	return &s, nil
}

func (db *DB) CreateSchedule(ctx context.Context, s *BookingSchedule) error {
	if s.TripInfo == nil {
		s.TripInfo = json.RawMessage("{}")
	}
	return db.Pool.QueryRow(ctx,
		`INSERT INTO booking_schedules (
			origin_keyword, origin_area_id, origin_name,
			dest_keyword, dest_area_id, dest_name,
			travel_date, time_from, time_to,
			seat_type, seat_count, auto_book,
			passenger_name, passenger_phone, passenger_email,
			status, trip_info
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,'pending',$16)
		RETURNING id, created_at, updated_at`,
		s.OriginKeyword, s.OriginAreaID, s.OriginName,
		s.DestKeyword, s.DestAreaID, s.DestName,
		s.TravelDate, s.TimeFrom, s.TimeTo,
		s.SeatType, s.SeatCount, s.AutoBook,
		s.PassengerName, s.PassengerPhone, s.PassengerEmail,
		s.TripInfo,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (db *DB) ListSchedules(ctx context.Context, statusFilter string) ([]BookingSchedule, error) {
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules`
	args := []any{}

	if statusFilter != "" && statusFilter != "all" {
		switch statusFilter {
		case "active":
			query += ` WHERE status IN ('pending','searching','found','booking')`
		case "success":
			query += ` WHERE status = 'success'`
		case "failed":
			query += ` WHERE status IN ('failed','cancelled')`
		default:
			query += ` WHERE status = $1`
			args = append(args, statusFilter)
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
	var s Stats
	err := db.Pool.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status IN ('searching','found','booking')),
			COUNT(*) FILTER (WHERE status = 'success'),
			COUNT(*) FILTER (WHERE status = 'failed'),
			COUNT(*) FILTER (WHERE status = 'cancelled')
		FROM booking_schedules
	`).Scan(&s.Total, &s.Pending, &s.Searching, &s.Success, &s.Failed, &s.Cancelled)
	return &s, err
}

func (db *DB) GetPendingSchedules(ctx context.Context, maxRetries int) ([]BookingSchedule, error) {
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules
		WHERE status IN ('pending', 'searching')
		AND retry_count < $1
		AND travel_date >= CURRENT_DATE
		ORDER BY travel_date ASC, created_at ASC`

	rows, err := db.Pool.Query(ctx, query, maxRetries)
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

func (db *DB) UpdateScheduleSuccess(ctx context.Context, id, bookingID, bookingCode, seatName, routeName string, price int, departureTime *time.Time) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE booking_schedules SET
			status='success', booking_id=$1, booking_code=$2,
			seat_name=$3, route_name=$4, ticket_price=$5, departure_time=$6,
			updated_at=NOW()
		WHERE id=$7`,
		bookingID, bookingCode, seatName, routeName, price, departureTime, id)
	return err
}

func (db *DB) MarkWebhookSent(ctx context.Context, id string) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE booking_schedules SET webhook_sent=TRUE, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (db *DB) GetRecentSchedules(ctx context.Context, limit int) ([]BookingSchedule, error) {
	query := `SELECT ` + scheduleColumns + ` FROM booking_schedules ORDER BY updated_at DESC LIMIT $1`
	rows, err := db.Pool.Query(ctx, query, limit)
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
