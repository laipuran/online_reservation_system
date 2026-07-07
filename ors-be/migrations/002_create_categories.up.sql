CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    parent_id   BIGINT       REFERENCES categories(id),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
