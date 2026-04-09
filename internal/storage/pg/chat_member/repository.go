package chat_member

import (
	"context"
	"database/sql"
)

type ChatMemberRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ChatMemberRepository {
	return &ChatMemberRepository{db: db}
}

func (r *ChatMemberRepository) IsMember(ctx context.Context, chatId int64, userId int64) (bool, error) {
	var exists bool

	err := r.db.QueryRowContext(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM users_to_chats
            WHERE chat_id = $1 AND user_id = $2
        )
    `, chatId, userId).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
