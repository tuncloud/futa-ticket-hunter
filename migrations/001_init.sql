-- App settings (single-row table for global config & user auth)
CREATE TABLE IF NOT EXISTS app_settings (
    id INT PRIMARY KEY DEFAULT 1,
    refresh_token TEXT DEFAULT '',
    full_name VARCHAR(255) DEFAULT '',
    phone VARCHAR(20) DEFAULT '',
    email VARCHAR(255) DEFAULT '',
    webhook_url TEXT DEFAULT '',
    webhook_secret TEXT DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT single_row CHECK (id = 1)
);

INSERT INTO app_settings (id) VALUES (1) ON CONFLICT DO NOTHING;

-- Booking schedules
CREATE TABLE IF NOT EXISTS booking_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Route info
    origin_keyword VARCHAR(255) NOT NULL,
    origin_area_id VARCHAR(100) DEFAULT '',
    origin_name VARCHAR(255) DEFAULT '',
    dest_keyword VARCHAR(255) NOT NULL,
    dest_area_id VARCHAR(100) DEFAULT '',
    dest_name VARCHAR(255) DEFAULT '',

    -- Time preferences
    travel_date DATE NOT NULL,
    time_from VARCHAR(10) DEFAULT '00:00',
    time_to VARCHAR(10) DEFAULT '23:59',

    -- Seat preferences
    seat_type VARCHAR(20) DEFAULT 'any',
    seat_count INT NOT NULL DEFAULT 1,
    auto_book BOOLEAN NOT NULL DEFAULT TRUE,

    -- Passenger info
    passenger_name VARCHAR(255) NOT NULL,
    passenger_phone VARCHAR(20) NOT NULL,
    passenger_email VARCHAR(255) DEFAULT '',

    -- Status: pending | searching | found | booking | success | failed | cancelled
    status VARCHAR(20) NOT NULL DEFAULT 'pending',

    -- Booking result
    booking_id VARCHAR(100) DEFAULT '',
    booking_code VARCHAR(20) DEFAULT '',
    ticket_price INT DEFAULT 0,
    seat_name VARCHAR(50) DEFAULT '',
    route_name VARCHAR(255) DEFAULT '',
    departure_time TIMESTAMPTZ,
    trip_info JSONB DEFAULT '{}',

    -- Metadata
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT DEFAULT '',
    webhook_sent BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_booking_schedules_status ON booking_schedules(status);
CREATE INDEX IF NOT EXISTS idx_booking_schedules_travel_date ON booking_schedules(travel_date);
