package message

import (
	"context"
	"database/sql"
	"github.com/thxhix/chat/internal/domain/chat"
	"github.com/thxhix/chat/internal/security/uuid"
)

type MessageRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) AddMessage(ctx context.Context, chatId int64, userId int64, text string) (string, error) {
	messageId, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	result, err := r.db.ExecContext(ctx, addMessageQuery, messageId, chatId, userId, text)
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
