package chat

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/thxhix/chat/internal/apperror"
	"github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/middleware"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Handler struct {
	logger     logger.ILogger
	jwtManager jwt.IJWTManager

	msgService message.IMessageService
}

func NewHandler(l logger.ILogger, jm jwt.IJWTManager, ms message.IMessageService) *Handler {
	return &Handler{
		logger:     l,
		jwtManager: jm,
		msgService: ms,
	}
}

type RequestMeta struct {
	ChatID uuid.UUID
	UserID int64
}

func getRequestMeta(r *http.Request) (*RequestMeta, error) {
	chatIDStr := chi.URLParam(r, CHAT_ID_KEY)
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return nil, apperror.NewBadRequestError("Wrong ChatID format provided")
	}

	userId, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		return nil, apperror.NewUnauthorizedError("Can't get user ID from token")
	}
	return &RequestMeta{
		ChatID: chatID,
		UserID: userId,
	}, nil
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// TODO...
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	rm, err := getRequestMeta(r)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	var req SendMessageRequest
	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	messageId, err := h.msgService.SendMessage(r.Context(), rm.ChatID, rm.UserID, req.Text)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	resp := SendMessageResponse{
		MessageID: messageId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := easyjson.MarshalToWriter(&resp, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}
