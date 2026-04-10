package apperror

import (
	"context"
	"errors"
	"net/http"
)

func HTTPStatus(err *Error) int {
	switch err.Code {
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeConflict:
		return http.StatusConflict
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeNotFound:
		return http.StatusNotFound
	case CodeForbidden:
		return http.StatusForbidden
	case CodeTimeout:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}

func ErrorMapper(err error) *Error {
	if errors.Is(err, context.DeadlineExceeded) {
		return NewTimeoutError("Request Timeout")
	}

	if errors.Is(err, context.Canceled) {
		return nil
	}

	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}

	return NewInternalError(err)
}
