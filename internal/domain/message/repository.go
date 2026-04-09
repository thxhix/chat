package message

import (
	"context"
)

type IMessageRepository interface {
	AddMessage(ctx context.Context, chatId int64, userId int64, text string) (string, error)
}
