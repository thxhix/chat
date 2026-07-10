package message

import (
	"context"
	"github.com/thxhix/chat/internal/domain/chat"
	uuidManager "github.com/thxhix/chat/internal/security/uuid"
	"github.com/thxhix/chat/internal/storage/pg/core/tx_manager"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
)

type IMessageService interface {
	GetChatMessages(ctx context.Context, chatId int64, userId int64, limit int, c *cursor.Cursor) (*core.Paged[MessageModel], error)
	SendMessage(ctx context.Context, chatId int64, userId int64, text string) (int64, error)
}

type MessageService struct {
	uuidManager uuidManager.IUUIDManager
	chatService chat.IChatService

	txManager   tx_manager.ITXManager
	messageRepo IMessageRepository
	chatRepo    chat.IChatRepository
}

func NewService(cs chat.IChatService, mr IMessageRepository, cr chat.IChatRepository, tx tx_manager.ITXManager, uuid uuidManager.IUUIDManager) *MessageService {
	return &MessageService{
		chatService: cs,
		messageRepo: mr,
		txManager:   tx,
		uuidManager: uuid,
	}
}

func (s *MessageService) GetChatMessages(ctx context.Context, chatId int64, userId int64, limit int, c *cursor.Cursor) (*core.Paged[MessageModel], error) {
	err := s.chatService.EnsureUserInChat(ctx, chatId, userId)
	if err != nil {
		return nil, err
	}

	messages, err := s.messageRepo.GetByChatID(ctx, chatId, limit, c)
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

func (s *MessageService) SendMessage(ctx context.Context, chatId int64, userId int64, text string) (int64, error) {
	var err error

	err = s.chatService.EnsureUserInChat(ctx, chatId, userId)
	if err != nil {
		return 0, err
	}

	mUUID, err := s.uuidManager.NewUUIDv7()
	if err != nil {
		return 0, err
	}

	var messageId int64
	err = s.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		messageId, err = s.messageRepo.AddMessage(ctx, mUUID, chatId, userId, text)
		if err != nil {
			return err
		}

		_, err = s.chatRepo.UpdateLastMsg(ctx, chatId, messageId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return messageId, nil
}
