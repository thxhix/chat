package token

import (
	"github.com/thxhix/chat/internal/apperror"
)

var (
	ErrInvalidCredentials           = apperror.NewValidationError("invalid credentials")
	ErrTokenAlreadyRotatedOrExpired = apperror.NewUnauthorizedError(`expired or invalid token`)
	ErrTokenDoesntExistsByJTI       = apperror.NewUnauthorizedError(`expired or invalid token`)
)
