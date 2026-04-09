package chat

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
)

type ChatRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetByUUID(ctx context.Context, chatId int64) (*chat.ChatModel, error) {
	query := `SELECT id, chat_id, type, title, created_at FROM chats WHERE chat_id = $1`

	var cr chat.ChatModel
	err := r.db.QueryRowContext(ctx, query, chatId).Scan(
		&cr.ID,
		&cr.ChatID,
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
	var id int64

	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM chats WHERE chat_id = $1`,
		chatUUID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
