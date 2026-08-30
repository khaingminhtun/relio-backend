-- +goose Up

CREATE TABLE memories (
    id BIGSERIAL PRIMARY KEY,

    relationship_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,

    title VARCHAR(255) NOT NULL,
    content TEXT,
    memory_date TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_memories_relationship
        FOREIGN KEY (relationship_id)
        REFERENCES relationships(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_memories_created_by
        FOREIGN KEY (created_by)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_memories_relationship_id
    ON memories(relationship_id);

CREATE INDEX idx_memories_created_by
    ON memories(created_by);

CREATE INDEX idx_memories_memory_date
    ON memories(relationship_id, memory_date DESC);

CREATE INDEX idx_memories_deleted_at
    ON memories(deleted_at);


-- +goose Down

DROP TABLE IF EXISTS memories;