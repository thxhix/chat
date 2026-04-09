package chat

import (
	"context"
	"github.com/google/uuid"
)

type IChatRepository interface {
	GetIDByUUID(ctx context.Context, chatUUID uuid.UUID) (int64, error)
	GetByUUID(ctx context.Context, chatId int64) (*ChatModel, error)
}
