package message

import (
	"context"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
	"github.com/thxhix/chat/internal/transport/http/core/result"
)

type IMessageService interface {
	GetChatMessages(ctx context.Context, chatId uuid.UUID, userId int64, limit int, c *cursor.Cursor) (*result.Paged[MessageModel], error)
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

func (s *MessageService) GetChatMessages(ctx context.Context, chatId uuid.UUID, userId int64, limit int, c *cursor.Cursor) (*result.Paged[MessageModel], error) {
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

	res := &result.Paged[MessageModel]{
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

	messageId, err := s.messageRepo.AddMessage(ctx, internalChatId, userId, text)
	if err != nil {
		return "", err
	}
	return messageId, nil
}
