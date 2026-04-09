package message

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
)

type IMessageService interface {
	SendMessage(ctx context.Context, chatId uuid.UUID, userId int64, text string) (string, error)
}

type MessageService struct {
	chatService chat.IChatService

	messageRepo IMessageRepository
}

func NewMessageService(cs chat.IChatService, mr IMessageRepository) *MessageService {
	return &MessageService{
		chatService: cs,
		messageRepo: mr,
	}
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

	messageId, err := s.messageRepo.AddMessage(ctx, internalChatId, userId, text)
	if err != nil {
		return "", err
	}
	return messageId, nil
}
