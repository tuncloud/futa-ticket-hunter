-- Add next_retry_at column to support concurrent workers with FOR UPDATE SKIP LOCKED.
-- A worker sets next_retry_at = NOW() + retry_delay when it claims a job, preventing
-- other workers (or other goroutines on the same worker) from picking up the same job
-- while it is being processed. If a worker crashes the job becomes claimable again
-- automatically once next_retry_at passes.
ALTER TABLE booking_schedules ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_booking_schedules_claim
    ON booking_schedules (status, travel_date, next_retry_at)
    WHERE status IN ('pending', 'searching');
