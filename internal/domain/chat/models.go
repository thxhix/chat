package chat

import (
	"github.com/google/uuid"
	"time"
)

type ChatModel struct {
	ID        int64
	ChatID    uuid.UUID
	Type      int8
	Title     string
	CreatedAt time.Time
}

type Chat struct {
	ID     int64
	ChatID uuid.UUID
	Type   int8
	Title  string
}
