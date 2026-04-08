package http

import (
	"encoding/json"
	"errors"
	"github.com/thxhix/chat/internal/apperror"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/transport/http/core"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrCantWriteResponseBody = errors.New("can't write response body")
)

func WriteError(w http.ResponseWriter, log logger.ILogger, err error) {
	var appErr *apperror.Error

	if !errors.As(err, &appErr) {
		appErr = &apperror.Error{
			Code:    apperror.CodeInternal,
			Message: "internal error",
			Err:     err,
		}
	}

	status := apperror.HTTPStatus(appErr)

	if status >= 500 {
		if appErr.Err != nil {
			log.Error("request failed", zap.Error(appErr.Err))
		} else {
			log.Error("request failed", zap.String("error", appErr.Message))
		}
	} else {
		log.Warn("request failed", zap.String("message", appErr.Message), zap.Int("status", status))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(core.ErrorResponse{
		Code:      status,
		ErrorText: appErr.Message,
	})
	if err != nil {
		log.Error("failed to write response", zap.Error(err))
	}
}
