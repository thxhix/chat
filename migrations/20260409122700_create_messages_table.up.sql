CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    message_id UUID NOT NULL UNIQUE,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL
);

CREATE INDEX idx_messages_chat_created_at ON messages(chat_id, created_at DESC);
CREATE INDEX idx_messages_not_deleted ON messages(chat_id, created_at DESC, id DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_messages_chat_cursor ON messages (chat_id, created_at DESC, id DESC);
