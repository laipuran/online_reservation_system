CREATE TABLE IF NOT EXISTS reservations (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id),
    service_id BIGINT      NOT NULL REFERENCES services(id),
    start_time TIMESTAMPTZ NOT NULL,
    end_time   TIMESTAMPTZ NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending'
                           CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled', 'rejected')),
    note       TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT reservations_time_range_check CHECK (end_time > start_time),
    CONSTRAINT reservations_service_start_time_unique UNIQUE (service_id, start_time)
);

CREATE INDEX IF NOT EXISTS idx_reservations_user_id ON reservations (user_id);
