CREATE TABLE IF NOT EXISTS notifications (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    title      VARCHAR(200) NOT NULL,
    content    TEXT         NOT NULL,
    type       VARCHAR(30)  NOT NULL
                            CHECK (type IN (
                                'reservation_confirmed',
                                'reservation_cancelled',
                                'reservation_reminder',
                                'review_received',
                                'system'
                            )),
    is_read    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_is_read ON notifications (user_id, is_read);
