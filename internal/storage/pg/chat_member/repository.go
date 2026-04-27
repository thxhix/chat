package chat_member

import (
	"context"
	"github.com/lib/pq"
	"github.com/thxhix/chat/internal/storage/pg/core"
)

type ChatMemberRepository struct {
	*core.BaseRepository
}

func NewRepository(base *core.BaseRepository) *ChatMemberRepository {
	return &ChatMemberRepository{
		BaseRepository: base,
	}
}

func (r *ChatMemberRepository) IsMember(ctx context.Context, chatId int64, userId int64) (bool, error) {
	executor := r.GetExecutor(ctx)

	var exists bool

	err := executor.QueryRowContext(ctx, queryIsMember, chatId, userId).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *ChatMemberRepository) Link(ctx context.Context, chatId int64, userId int64) error {
	executor := r.GetExecutor(ctx)
	_, err := executor.ExecContext(ctx, queryLink, chatId, userId)
	return err
}

func (r *ChatMemberRepository) LinkBatch(ctx context.Context, chatId int64, userList []int64) error {
	executor := r.GetExecutor(ctx)
	_, err := executor.ExecContext(ctx, queryLinkBatch, chatId, pq.Array(userList))
	return err
}
