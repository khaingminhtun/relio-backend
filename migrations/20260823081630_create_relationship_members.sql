-- +goose Up
CREATE TABLE relationship_members (
                                      id BIGSERIAL PRIMARY KEY,

                                      relationship_id BIGINT NOT NULL,
                                      user_id BIGINT NOT NULL,

                                      role VARCHAR(50) NOT NULL,
                                      status VARCHAR(50) NOT NULL,

                                      joined_at TIMESTAMPTZ NULL,

                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      deleted_at TIMESTAMPTZ NULL,

                                      CONSTRAINT fk_relationship_members_relationship
                                          FOREIGN KEY (relationship_id)
                                              REFERENCES relationships(id)
                                              ON DELETE RESTRICT,

                                      CONSTRAINT fk_relationship_members_user
                                          FOREIGN KEY (user_id)
                                              REFERENCES users(id)
                                              ON DELETE RESTRICT
);

CREATE INDEX idx_relationship_members_relationship_id
    ON relationship_members(relationship_id);

CREATE INDEX idx_relationship_members_user_id
    ON relationship_members(user_id);

CREATE INDEX idx_relationship_members_deleted_at
    ON relationship_members(deleted_at);

CREATE INDEX idx_relationship_members_status
    ON relationship_members(status);

CREATE UNIQUE INDEX idx_relationship_members_active_unique
    ON relationship_members(relationship_id, user_id)
    WHERE deleted_at IS NULL;


-- +goose Down
DROP TABLE relationship_members;