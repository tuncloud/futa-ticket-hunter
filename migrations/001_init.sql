-- Booking schedules
CREATE TABLE IF NOT EXISTS booking_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Route info
    origin_area_id VARCHAR(100) NOT NULL,
    origin_name VARCHAR(255) NOT NULL,
    dest_area_id VARCHAR(100) NOT NULL,
    dest_name VARCHAR(255) NOT NULL,

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
    route_name VARCHAR(500) DEFAULT '',
    departure_time TIMESTAMPTZ,

    -- Metadata
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_booking_schedules_status ON booking_schedules(status);
CREATE INDEX IF NOT EXISTS idx_booking_schedules_travel_date ON booking_schedules(travel_date);
CREATE INDEX IF NOT EXISTS idx_booking_schedules_email ON booking_schedules(passenger_email);
