package message

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -all dto.go

type SendMessageRequest struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type SendMessageResponse struct {
	MessageID int64 `json:"message_id"`
}

type GetMessagesRequest struct {
	ChatId int64  `json:"chat_id"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

type GetMessagesRecord struct {
	MessageID uuid.UUID `json:"message_id"`
	Text      string    `json:"text"`
	From      int64     `json:"from"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetMessagesResponse struct {
	Messages   []GetMessagesRecord `json:"messages"`
	NextCursor *string             `json:"next_cursor"`
}
