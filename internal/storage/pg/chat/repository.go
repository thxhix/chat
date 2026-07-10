package chat

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
	"github.com/thxhix/chat/internal/storage/pg/core"
)

type ChatRepository struct {
	*core.BaseRepository
}

func NewRepository(base *core.BaseRepository) *ChatRepository {
	return &ChatRepository{
		BaseRepository: base,
	}
}

func (r *ChatRepository) CreateChat(ctx context.Context, idempotencyKey uuid.UUID, cType int8) (*chat.CreateChatResult, bool, error) {
	executor := r.GetExecutor(ctx)

	res := &chat.CreateChatResult{}
	var isCreated bool

	err := executor.QueryRowContext(ctx, queryCreateChat, idempotencyKey, cType, 0).Scan(
		&res.ID,
		&res.IdempotencyKey,
		&isCreated,
	)
	if err != nil {
		return res, false, err
	}

	return res, isCreated, nil
}

func (r *ChatRepository) GetUserChats(ctx context.Context, userId int64, limit int64) ([]*chat.ChatModel, error) {
	executor := r.GetExecutor(ctx)

	rows, err := executor.QueryContext(ctx, queryGetUserChats, userId, limit)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	chats := make([]*chat.ChatModel, 0, limit)
	for rows.Next() {
		row := &chat.ChatModel{}
		err = rows.Scan(
			&row.ID,
			&row.IdempotencyKey,
			&row.Type,
			&row.Title,
			&row.CreatedAt,
			&row.LastMsgID,
			&row.LastMsgText,
		)
		if err != nil {
			return nil, err
		}

		chats = append(chats, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepository) UpdateLastMsg(ctx context.Context, chatId int64, msgId int64) (bool, error) {
	executor := r.GetExecutor(ctx)

	_, err := executor.ExecContext(ctx, queryUpdateLastMsg, chatId, msgId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *ChatRepository) GetByUUID(ctx context.Context, chatId int64) (*chat.ChatModel, error) {
	executor := r.GetExecutor(ctx)

	var cr chat.ChatModel
	err := executor.QueryRowContext(ctx, queryGetByUUID, chatId).Scan(
		&cr.ID,
		&cr.IdempotencyKey,
		&cr.Type,
		&cr.Title,
		&cr.CreatedAt,
	)
	if err != nil {
		return &chat.ChatModel{}, err
	}

	return &cr, nil
}

func (r *ChatRepository) GetIDByUUID(ctx context.Context, chatUUID uuid.UUID) (int64, error) {
	executor := r.GetExecutor(ctx)

	var id int64

	err := executor.QueryRowContext(ctx, queryGetIDByUUID, chatUUID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
