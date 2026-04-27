package chat

import "github.com/google/uuid"

//go:generate easyjson -all dto.go

type CreateChatRequest struct {
	ChatType int8     `json:"chat_type"`
	Members  []string `json:"members"`
}

type CreateChatResponse struct {
	ChatId   int64     `json:"chat_id"`
	ChatUUID uuid.UUID `json:"chat_uuid"`
}
