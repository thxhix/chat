package chat

import (
	"context"
	"github.com/google/uuid"
)

type IChatRepository interface {
	CreateChat(ctx context.Context, idempotencyKey uuid.UUID, cType int8) (*CreateChatResult, bool, error)
	GetUserChats(ctx context.Context, userId int64, limit int64) ([]*ChatModel, error)

	UpdateLastMsg(ctx context.Context, chatId int64, msgId int64) (bool, error)

	GetIDByUUID(ctx context.Context, chatUUID uuid.UUID) (int64, error)
	GetByUUID(ctx context.Context, chatId int64) (*ChatModel, error)
}
