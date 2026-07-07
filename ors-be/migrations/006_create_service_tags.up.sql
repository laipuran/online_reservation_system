CREATE TABLE IF NOT EXISTS service_tags (
    service_id BIGINT NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    tag_id     BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (service_id, tag_id)
);

CREATE INDEX idx_service_tags_tag_id ON service_tags (tag_id);
