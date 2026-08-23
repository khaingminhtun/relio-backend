-- +goose Up

CREATE TABLE relationships (
                               id BIGSERIAL PRIMARY KEY,

                               name VARCHAR(150) NOT NULL,
                               type VARCHAR(50) NOT NULL,
                               custom_type VARCHAR(100) NULL,
                               description TEXT NULL,
                               start_date TIMESTAMPTZ NULL,
                               timezone VARCHAR(100) NOT NULL,

                               created_by BIGINT NOT NULL,

                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               deleted_at TIMESTAMPTZ NULL,

                               CONSTRAINT fk_relationships_created_by
                                   FOREIGN KEY (created_by)
                                       REFERENCES users(id)
                                       ON DELETE RESTRICT
);

CREATE INDEX idx_relationships_created_by
    ON relationships(created_by);

CREATE INDEX idx_relationships_deleted_at
    ON relationships(deleted_at);

CREATE INDEX idx_relationships_type
    ON relationships(type);


-- +goose Down

DROP TABLE relationships;