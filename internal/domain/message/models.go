package message

import (
	"github.com/google/uuid"
	"time"
)

type MessageModel struct {
	ID        int64
	MessageID uuid.UUID
	ChatID    int64
	UserID    int64
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
