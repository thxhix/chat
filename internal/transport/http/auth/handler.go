package auth

import (
	"github.com/mailru/easyjson"
	"github.com/thxhix/chat/internal/domain/auth"
	"github.com/thxhix/chat/internal/domain/token"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/transport/http/core"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type Handler struct {
	logger      logger.ILogger
	authService auth.IAuthService
}

func NewHandler(l logger.ILogger, as auth.IAuthService) *Handler {
	return &Handler{
		logger:      l,
		authService: as,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	reqObj := RegisterRequest{}
	err = easyjson.Unmarshal(body, &reqObj)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	userId, accessToken, refreshToken, err := h.authService.Register(r.Context(), reqObj.Login, reqObj.Password)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	respObj := RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusCreated)

	if _, err := easyjson.MarshalToWriter(&respObj, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	reqObj := LoginRequest{}
	err = easyjson.Unmarshal(body, &reqObj)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	userId, accessToken, refreshToken, err := h.authService.Login(r.Context(), reqObj.Login, reqObj.Password)

	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	respObj := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	if _, err := easyjson.MarshalToWriter(&respObj, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	reqObj := RefreshRequest{}
	err = easyjson.Unmarshal(body, &reqObj)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	if err := token.ValidateRefreshToken(reqObj.RefreshToken); err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	userId, accessToken, refreshToken, err := h.authService.Refresh(r.Context(), reqObj.RefreshToken)
	if err != nil {
		core.WriteError(w, h.logger, err)
		return
	}

	respObj := RefreshedTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)

	if _, err := easyjson.MarshalToWriter(&respObj, w); err != nil {
		h.logger.Error(core.ErrCantWriteResponseBody.Error(), zap.Error(err))
		return
	}
}
