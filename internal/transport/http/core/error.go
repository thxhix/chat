package core

import (
	"encoding/json"
	"errors"
	"github.com/thxhix/chat/internal/apperror"
	"github.com/thxhix/chat/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

var (
	ErrCantWriteResponseBody = errors.New("can't write response body")
)

func WriteError(w http.ResponseWriter, log logger.ILogger, err error) {
	appErr := apperror.ErrorMapper(err)
	if appErr == nil {
		return
	}
	status := apperror.HTTPStatus(appErr)

	writeLogger(log, status, appErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err = json.NewEncoder(w).Encode(ErrorResponse{
		Code:      status,
		ErrorText: appErr.Message,
	})
	if err != nil {
		log.Error("failed to write response", zap.Error(err))
	}
}

func writeLogger(logger logger.ILogger, statusCode int, appErr *apperror.Error) {
	if statusCode >= 500 {
		if appErr.Err != nil {
			logger.Error("INTERNAL Error", zap.Int("status", statusCode), zap.Error(appErr.Err))
		} else {
			logger.Error("INTERNAL Error", zap.Int("status", statusCode), zap.String("message", appErr.Message))
		}
	} else {
		logger.Warn("Request Error", zap.Int("status", statusCode), zap.String("message", appErr.Message))
	}
}
