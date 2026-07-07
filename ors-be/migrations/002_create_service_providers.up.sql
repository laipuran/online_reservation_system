CREATE TABLE IF NOT EXISTS service_providers (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    business_name VARCHAR(200) NOT NULL,
    description   TEXT,
    address       VARCHAR(500),
    phone         VARCHAR(20),
    email         VARCHAR(255),
    logo_url      VARCHAR(500),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_providers_user_id ON service_providers (user_id);
