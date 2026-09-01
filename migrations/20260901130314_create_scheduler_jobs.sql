-- +goose Up

CREATE TABLE scheduled_jobs (
    id BIGSERIAL PRIMARY KEY,

    relationship_id BIGINT NOT NULL,
    event_id BIGINT NOT NULL,
    reminder_id BIGINT NOT NULL,

    type VARCHAR(50) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',

    scheduled_at TIMESTAMPTZ NOT NULL,

    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,

    last_error TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_scheduled_jobs_relationship
        FOREIGN KEY (relationship_id)
        REFERENCES relationships(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_scheduled_jobs_event
        FOREIGN KEY (event_id)
        REFERENCES events(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_scheduled_jobs_reminder
        FOREIGN KEY (reminder_id)
        REFERENCES reminders(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_scheduled_jobs_due
    ON scheduled_jobs (status, scheduled_at);

CREATE INDEX idx_scheduled_jobs_event_id
    ON scheduled_jobs (event_id);

CREATE INDEX idx_scheduled_jobs_reminder_id
    ON scheduled_jobs (reminder_id);

CREATE INDEX idx_scheduled_jobs_relationship_id
    ON scheduled_jobs (relationship_id);


-- +goose Down

DROP TABLE IF EXISTS scheduled_jobs;