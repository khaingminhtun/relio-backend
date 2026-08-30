-- +goose Up

CREATE TABLE notes (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,

    mood VARCHAR(50),

    is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_notes_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_notes_user_id
    ON notes(user_id);

CREATE INDEX idx_notes_user_created_at
    ON notes(user_id, created_at DESC);

CREATE INDEX idx_notes_user_pinned
    ON notes(user_id, is_pinned)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_notes_user_archived
    ON notes(user_id, is_archived)
    WHERE deleted_at IS NULL;


-- +goose Down

DROP TABLE IF EXISTS notes;