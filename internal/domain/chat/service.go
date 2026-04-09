package chat

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat_member"
)

type IChatService interface {
	GetInternalID(ctx context.Context, chatUUID uuid.UUID) (int64, error)
	EnsureUserInChat(ctx context.Context, chatId int64, userId int64) error
}

type ChatService struct {
	chatRepo       IChatRepository
	chatMemberRepo chat_member.IChatMemberRepository
}

// NewService constructs a new AuthService with given dependencies.
func NewService(cr IChatRepository, cmr chat_member.IChatMemberRepository) IChatService {
	return &ChatService{
		chatRepo:       cr,
		chatMemberRepo: cmr,
	}
}

func (s *ChatService) GetInternalID(ctx context.Context, chatUUID uuid.UUID) (int64, error) {
	id, err := s.chatRepo.GetIDByUUID(ctx, chatUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (s *ChatService) EnsureUserInChat(ctx context.Context, chatId int64, userId int64) error {
	ok, err := s.chatMemberRepo.IsMember(ctx, chatId, userId)
	if err != nil {
		return err
	}

	if !ok {
		return ErrForbidden
	}

	return nil
}
