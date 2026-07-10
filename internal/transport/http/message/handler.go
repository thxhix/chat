package message

import (
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/thxhix/chat/internal/apperror"
	"github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/core/cursor"
	"github.com/thxhix/chat/internal/transport/http/core/limit"
	"github.com/thxhix/chat/internal/transport/http/middleware"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	logger     logger.ILogger
	msgService message.IMessageService
}

func NewHandler(l logger.ILogger, ms message.IMessageService) *Handler {
	return &Handler{
		logger:     l,
		msgService: ms,
	}
}

type RequestMeta struct {
	ChatID int64
	UserID int64
}

func getRequestMeta(r *http.Request) (*RequestMeta, error) {
	chatIdStr := chi.URLParam(r, core.ChatIdPrefix)
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid chat_id")
	}

	userId, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		return nil, apperror.NewUnauthorizedError("Can't get user ID from token")
	}
	return &RequestMeta{
		ChatID: chatId,
		UserID: userId,
	}, nil
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	rm, err := getRequestMeta(r)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	reqLimit := limit.GetFromRequest(r)
	reqCursorStr := cursor.GetFromRequest(r)

	reqCursor, err := cursor.DecodeCursor(reqCursorStr)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	result, err := h.msgService.GetChatMessages(r.Context(), rm.ChatID, rm.UserID, reqLimit, reqCursor)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	var nextCursor *string
	if result.HasMore && result.LastItem != nil {
		c := cursor.Cursor{
			ID:        result.LastItem.ID,
			CreatedAt: result.LastItem.CreatedAt,
		}
		encoded, _ := c.EncodeCursor()
		nextCursor = &encoded
	}

	mList := make([]GetMessagesRecord, 0, len(result.Items))
	for _, m := range result.Items {
		mList = append(mList, ToMessageResponse(m))
	}

	res := GetMessagesResponse{
		Messages:   mList,
		NextCursor: nextCursor,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := easyjson.MarshalToWriter(&res, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	rm, err := getRequestMeta(r)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var req SendMessageRequest
	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	messageId, err := h.msgService.SendMessage(r.Context(), req.ChatId, rm.UserID, req.Text)
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
