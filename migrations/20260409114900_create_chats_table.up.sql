CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    chat_id UUID NOT NULL UNIQUE,
    type SMALLINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chats_type ON chats(type);
