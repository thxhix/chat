package chat_member

import (
	"context"
)

type IChatMemberRepository interface {
	IsMember(ctx context.Context, chatId int64, userId int64) (bool, error)
}
