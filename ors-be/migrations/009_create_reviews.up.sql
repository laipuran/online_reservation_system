CREATE TABLE IF NOT EXISTS reviews (
    id             BIGSERIAL   PRIMARY KEY,
    reservation_id BIGINT      NOT NULL UNIQUE REFERENCES reservations(id),
    user_id        BIGINT      NOT NULL REFERENCES users(id),
    service_id     BIGINT      NOT NULL REFERENCES services(id),
    rating         SMALLINT    NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment        TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_service_id ON reviews (service_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews (user_id);
