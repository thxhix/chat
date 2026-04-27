package message

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
	"github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/storage/pg/core"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
)

type MessageRepository struct {
	*core.BaseRepository
}

func NewRepository(base *core.BaseRepository) *MessageRepository {
	return &MessageRepository{
		BaseRepository: base,
	}
}

func (r *MessageRepository) GetByChatID(ctx context.Context, chatId int64, limit int, c *cursor.Cursor) ([]message.MessageModel, error) {
	if c == nil {
		return r.getByChatIDFirstPage(ctx, chatId, limit)
	}
	return r.getByChatIDCursor(ctx, chatId, limit, c)
}

func (r *MessageRepository) getByChatIDFirstPage(ctx context.Context, chatId int64, limit int) ([]message.MessageModel, error) {
	rows, err := r.BaseRepository.DB.QueryContext(ctx, queryGetByChatMessagesFirstPage, chatId, limit)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	messages := make([]message.MessageModel, 0, limit)
	for rows.Next() {
		row := message.MessageModel{}
		err = rows.Scan(
			&row.ID,
			&row.MessageID,
			&row.ChatID,
			&row.UserID,
			&row.Text,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) getByChatIDCursor(ctx context.Context, chatId int64, limit int, c *cursor.Cursor) ([]message.MessageModel, error) {
	rows, err := r.BaseRepository.DB.QueryContext(ctx, queryGetByChatMessagesWithCursor, chatId, c.CreatedAt, c.ID, limit)
	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	messages := make([]message.MessageModel, 0, limit)
	for rows.Next() {
		row := message.MessageModel{}
		err = rows.Scan(
			&row.ID,
			&row.MessageID,
			&row.ChatID,
			&row.UserID,
			&row.Text,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) AddMessage(ctx context.Context, messageId uuid.UUID, chatId int64, userId int64, text string) (string, error) {
	result, err := r.BaseRepository.DB.ExecContext(ctx, addMessageQuery, messageId, chatId, userId, text)
	if err != nil {
		return "", err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if affected == 0 {
		return "", chat.ErrForbidden
	}

	return messageId.String(), nil
}
