package chat

import (
	"github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"net/http"
)

type Handler struct {
	logger logger.ILogger

	msgService message.IMessageService
}

func NewHandler(l logger.ILogger, ms message.IMessageService) *Handler {
	return &Handler{
		logger:     l,
		msgService: ms,
	}
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	return
}
