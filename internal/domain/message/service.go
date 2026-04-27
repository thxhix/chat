package message

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
	uuidManager "github.com/thxhix/chat/internal/security/uuid"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
)

type IMessageService interface {
	GetChatMessages(ctx context.Context, chatId uuid.UUID, userId int64, limit int, c *cursor.Cursor) (*core.Paged[MessageModel], error)
	SendMessage(ctx context.Context, chatId uuid.UUID, userId int64, text string) (string, error)
}

type MessageService struct {
	uuidManager uuidManager.IUUIDManager
	chatService chat.IChatService

	messageRepo IMessageRepository
}

func NewMessageService(cs chat.IChatService, mr IMessageRepository, uuid uuidManager.IUUIDManager) *MessageService {
	return &MessageService{
		chatService: cs,
		messageRepo: mr,
		uuidManager: uuid,
	}
}

func (s *MessageService) GetChatMessages(ctx context.Context, chatId uuid.UUID, userId int64, limit int, c *cursor.Cursor) (*core.Paged[MessageModel], error) {
	internalChatId, err := s.chatService.GetInternalID(ctx, chatId)
	if err != nil {
		return nil, err
	}

	err = s.chatService.EnsureUserInChat(ctx, internalChatId, userId)
	if err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.GetByChatID(ctx, internalChatId, limit, c)
	if err != nil {
		return nil, err
	}

	res := &core.Paged[MessageModel]{
		Items: messages,
	}

	if len(messages) > limit {

		res.HasMore = true
		res.Items = messages[:limit]
		res.LastItem = &res.Items[len(res.Items)-1]
	}

	return res, nil
}

func (s *MessageService) SendMessage(ctx context.Context, chatId uuid.UUID, userId int64, text string) (string, error) {
	internalChatId, err := s.chatService.GetInternalID(ctx, chatId)
	if err != nil {
		return "", err
	}

	err = s.chatService.EnsureUserInChat(ctx, internalChatId, userId)
	if err != nil {
		return "", err
	}

	mUUID, err := s.uuidManager.NewUUIDv7()
	if err != nil {
		return "", err
	}

	messageId, err := s.messageRepo.AddMessage(ctx, mUUID, internalChatId, userId, text)
	if err != nil {
		return "", err
	}
	return messageId, nil
}
