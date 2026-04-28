package chat

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat_member"
	uuidManager "github.com/thxhix/chat/internal/security/uuid"
	"github.com/thxhix/chat/internal/storage/pg/core/tx_manager"
)

type IChatService interface {
	CreateChat(ctx context.Context, cType int8, members []int64) (*CreateChatResult, bool, error)
	GetUserChats(ctx context.Context, userId int64) (*GetChatsResult, error)

	GetInternalID(ctx context.Context, chatUUID uuid.UUID) (int64, error)
	EnsureUserInChat(ctx context.Context, chatId int64, userId int64) error
}

type ChatService struct {
	chatRepo       IChatRepository
	chatMemberRepo chat_member.IChatMemberRepository
	txManager      tx_manager.ITXManager
	uuidManager    uuidManager.IUUIDManager
}

// NewService constructs a new AuthService with given dependencies.
func NewService(cr IChatRepository, cmr chat_member.IChatMemberRepository, tx tx_manager.ITXManager, uuid uuidManager.IUUIDManager) IChatService {
	return &ChatService{
		chatRepo:       cr,
		chatMemberRepo: cmr,

		txManager:   tx,
		uuidManager: uuid,
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

func (s *ChatService) CreateChat(ctx context.Context, cType int8, members []int64) (*CreateChatResult, bool, error) {
	var err error

	var chat *CreateChatResult
	var idempotencyKey uuid.UUID
	var isCreated bool

	if cType == 1 || cType == 3 {
		idempotencyKey = s.uuidManager.NewUUIDv5(GenerateIdempotencyStr(cType, members))
	} else {
		idempotencyKey, err = s.uuidManager.NewUUIDv7()
		if err != nil {
			return chat, false, err
		}
	}

	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		chat, isCreated, err = s.chatRepo.CreateChat(txCtx, idempotencyKey, cType)
		if err != nil {
			return err
		}

		err = s.chatMemberRepo.LinkBatch(txCtx, chat.ID, members)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return chat, false, err
	}

	return chat, isCreated, nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userId int64) (*GetChatsResult, error) {
	chats, err := s.chatRepo.GetUserChats(ctx, userId, 25)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(chats))
	for i, chat := range chats {
		ids[i] = chat.ID
	}

	membersMap, err := s.chatMemberRepo.GetMembersForChats(ctx, ids, userId)
	if err != nil {
		return nil, err
	}

	res := make([]*Chat, 0, len(chats))

	for _, c := range chats {
		chatMembers := membersMap[c.ID]

		row := &Chat{
			ID:             c.ID,
			IdempotencyKey: c.IdempotencyKey,
			Type:           c.Type,
			CreatedAt:      c.CreatedAt,
		}

		if c.Type != 3 {
			if len(chatMembers) > 0 {
				row.Participants = chatMembers[0]
			}
		}

		res = append(res, row)
	}

	return &GetChatsResult{
		Chats: res,
	}, nil
}
