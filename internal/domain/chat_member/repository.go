package chat_member

import (
	"context"
)

type IChatMemberRepository interface {
	Link(ctx context.Context, chatId int64, userId int64) error
	LinkBatch(ctx context.Context, chatId int64, userList []int64) error

	IsMember(ctx context.Context, chatId int64, userId int64) (bool, error)
	GetMembersForChats(ctx context.Context, chatIDs []int64, userId int64) (map[int64][]*Member, error)
}
