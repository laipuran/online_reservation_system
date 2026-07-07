CREATE TABLE IF NOT EXISTS services (
    id               BIGSERIAL      PRIMARY KEY,
    provider_id      BIGINT         NOT NULL REFERENCES service_providers(id) ON DELETE CASCADE,
    category_id      BIGINT         NOT NULL REFERENCES categories(id),
    title            VARCHAR(200)   NOT NULL,
    description      TEXT,
    price            NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    duration_minutes INTEGER        NOT NULL CHECK (duration_minutes > 0),
    image_url        VARCHAR(500),
    status           VARCHAR(20)    NOT NULL DEFAULT 'active'
                                    CHECK (status IN ('active', 'inactive', 'pending', 'rejected')),
    avg_rating       REAL           NOT NULL DEFAULT 0 CHECK (avg_rating >= 0 AND avg_rating <= 5),
    review_count     INTEGER        NOT NULL DEFAULT 0 CHECK (review_count >= 0),
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_services_provider_id ON services (provider_id);
CREATE INDEX idx_services_category_id ON services (category_id);
CREATE INDEX idx_services_status ON services (status);
CREATE INDEX idx_services_created_at ON services (created_at);
