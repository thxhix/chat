package message

import (
	"github.com/thxhix/chat/internal/domain/message"
)

func ToMessageResponse(m message.MessageModel) GetMessagesRecord {
	return GetMessagesRecord{
		MessageID: m.MessageID,
		Text:      m.Text,
		From:      m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
