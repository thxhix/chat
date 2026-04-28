package chat_member

import (
	"context"
	"github.com/lib/pq"
	"github.com/thxhix/chat/internal/domain/chat_member"
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

func (r *ChatMemberRepository) GetMembersForChats(ctx context.Context, chatIDs []int64, userId int64) (map[int64][]*chat_member.Member, error) {
	executor := r.GetExecutor(ctx)

	rows, err := executor.QueryContext(ctx, queryGetMembers, pq.Array(chatIDs), userId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	res := make(map[int64][]*chat_member.Member)
	for rows.Next() {
		m := &chat_member.Member{}

		err = rows.Scan(
			&m.ChatID,
			&m.UserID,
			&m.UserLogin,
		)
		if err != nil {
			return nil, err
		}

		res[m.ChatID] = append(res[m.ChatID], m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
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
