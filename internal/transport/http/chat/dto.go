package chat

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -all dto.go

type CreateChatRequest struct {
	ChatType int8     `json:"chat_type"`
	Members  []string `json:"members"`
}

type CreateChatResponse struct {
	ChatId   int64     `json:"chat_id"`
	ChatUUID uuid.UUID `json:"chat_uuid"`
}

type Chat struct {
	ID              int64          `json:"id"`
	UUID            uuid.UUID      `json:"uuid"`
	Title           *string        `json:"title"`
	Type            int8           `json:"type"`
	CreatedAt       time.Time      `json:"created_at"`
	Target          *ChatUserShort `json:"target,omitempty"`
	LastMessageText *string        `json:"last_message_text"`
}

type ChatUserShort struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type GetChatsResponse struct {
	Chats []*Chat `json:"chats"`
}
