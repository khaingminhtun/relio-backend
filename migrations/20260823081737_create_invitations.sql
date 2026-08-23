-- +goose Up

CREATE TABLE invitations (
                             id BIGSERIAL PRIMARY KEY,

                             relationship_id BIGINT NOT NULL,
                             invited_by BIGINT NOT NULL,

                             email VARCHAR(255) NOT NULL,
                             token_hash VARCHAR(255) NOT NULL UNIQUE,

                             status VARCHAR(50) NOT NULL,
                             expires_at TIMESTAMPTZ NOT NULL,
                             accepted_at TIMESTAMPTZ NULL,

                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                             CONSTRAINT fk_invitations_relationship
                                 FOREIGN KEY (relationship_id)
                                     REFERENCES relationships(id)
                                     ON DELETE RESTRICT,

                             CONSTRAINT fk_invitations_invited_by
                                 FOREIGN KEY (invited_by)
                                     REFERENCES users(id)
                                     ON DELETE RESTRICT
);

CREATE INDEX idx_invitations_relationship_id
    ON invitations(relationship_id);

CREATE INDEX idx_invitations_invited_by
    ON invitations(invited_by);

CREATE INDEX idx_invitations_email
    ON invitations(email);

CREATE INDEX idx_invitations_status
    ON invitations(status);

CREATE INDEX idx_invitations_expires_at
    ON invitations(expires_at);


-- +goose Down

DROP TABLE invitations;