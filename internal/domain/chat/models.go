package chat

import (
	"github.com/google/uuid"
	"github.com/thxhix/chat/internal/domain/chat_member"
	"time"
)

type ChatModel struct {
	ID             int64
	IdempotencyKey uuid.UUID
	Type           int8
	Title          *string
	LastMsgID      *int64
	LastMsgText    *string
	CreatedAt      time.Time
}

type Chat struct {
	ID             int64
	IdempotencyKey uuid.UUID
	Type           int8
	Title          *string
	CreatedAt      time.Time
	Participants   *chat_member.Member
	LastMsgText    *string
}

type CreateChatResult struct {
	ID             int64
	IdempotencyKey uuid.UUID
}

type GetChatsResult struct {
	Chats []*Chat
}
