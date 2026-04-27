package message

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
)

type IMessageRepository interface {
	GetByChatID(ctx context.Context, chatId int64, limit int, c *cursor.Cursor) ([]MessageModel, error)
	AddMessage(ctx context.Context, messageId uuid.UUID, chatId int64, userId int64, text string) (string, error)
}
