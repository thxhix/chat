CREATE TABLE users_to_chats (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role SMALLINT NOT NULL DEFAULT 0,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    left_at TIMESTAMPTZ NULL DEFAULT NULL,
    
    UNIQUE (user_id, chat_id)
);

CREATE INDEX idx_utc_user_chat ON users_to_chats(user_id, chat_id);
CREATE INDEX idx_utc_chat_user ON users_to_chats(chat_id, user_id);
