package chat

import (
	"github.com/mailru/easyjson"
	"github.com/thxhix/chat/internal/apperror"
	"github.com/thxhix/chat/internal/domain/chat"
	"github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/middleware"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	logger logger.ILogger

	msgService  message.IMessageService
	chatService chat.IChatService
}

func NewHandler(l logger.ILogger, ms message.IMessageService, cs chat.IChatService) *Handler {
	return &Handler{
		logger:      l,
		msgService:  ms,
		chatService: cs,
	}
}

func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	var req CreateChatRequest
	err = easyjson.Unmarshal(body, &req)
	if err != nil {
		core.WriteError(w, h.logger, apperror.NewBadRequestError("Invalid request body"))
		return
	}

	err = chat.ValidateCreateRequest(req.ChatType, req.Members)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	userId, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		core.WriteError(w, h.logger, apperror.NewUnauthorizedError("Can't get user ID from token"))
		return
	}

	memberIDs := make([]int64, len(req.Members))
	for i, s := range req.Members {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			core.WriteError(w, h.logger, apperror.NewBadRequestError("Cant't parse new chat member IDs"))
			return
		}
		memberIDs[i] = id
	}
	memberIDs = append(memberIDs, userId)

	resChat, isCreated, err := h.chatService.CreateChat(r.Context(), req.ChatType, memberIDs)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	res := CreateChatResponse{
		ChatId:   resChat.ID,
		ChatUUID: resChat.IdempotencyKey,
	}

	w.Header().Set("Content-Type", "application/json")
	if isCreated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if _, err := easyjson.MarshalToWriter(&res, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}
