package token

import (
	"github.com/thxhix/chat/internal/apperror"
)

var (
	ErrInvalidCredentials = apperror.NewValidationError("invalid credentials")
	ErrDeadToken          = apperror.NewUnauthorizedError(`expired or invalid token`)
	ErrMissingToken       = apperror.NewUnauthorizedError(`missing bearer token`)
)
