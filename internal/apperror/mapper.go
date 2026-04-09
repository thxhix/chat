package apperror

import (
	"errors"
	"net/http"
)

func HTTPStatus(err error) int {
	var appErr *Error

	if errors.As(err, &appErr) {
		switch appErr.Code {
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
		default:
			return http.StatusInternalServerError
		}
	}

	return http.StatusInternalServerError
}
