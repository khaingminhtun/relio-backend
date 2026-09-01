-- +goose Up

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,

    relationship_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    description TEXT,

    event_date DATE NOT NULL,
    event_time TIME,
    all_day BOOLEAN NOT NULL DEFAULT FALSE,

    created_by BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_events_relationship
        FOREIGN KEY (relationship_id)
        REFERENCES relationships(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_events_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_events_relationship_id
    ON events (relationship_id);

CREATE INDEX idx_events_event_date
    ON events (event_date);

CREATE INDEX idx_events_created_by
    ON events (created_by);

CREATE INDEX idx_events_deleted_at
    ON events (deleted_at);


-- +goose Down

DROP TABLE IF EXISTS events;