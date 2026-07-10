CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    idempotency_key uuid NOT NULL,
    type SMALLINT NOT NULL,
    title VARCHAR(255) NULL DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    last_msg_id BIGINT NOT NULL,

    UNIQUE (idempotency_key)
);

CREATE INDEX idx_chats_type ON chats(type);
CREATE INDEX idx_chats_idempotency_key ON chats(idempotency_key);
CREATE INDEX idx_chats_last_msg ON chats(last_msg_id);
