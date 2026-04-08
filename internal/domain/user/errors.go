package user

import (
	"github.com/thxhix/chat/internal/apperror"
)

var (
	ErrLoginTooShort    = apperror.NewValidationError("login is too short, minimum length – 3")
	ErrLoginTooLong     = apperror.NewValidationError("login is too long, maximum length – 64")
	ErrPasswordTooShort = apperror.NewValidationError("password too short, minimum length – 8")
	ErrPasswordTooLong  = apperror.NewValidationError("password too long, maximum length – 50")
)

var (
	ErrDuplicateLogin = apperror.NewConflictError("login is already taken")
	ErrUserNotFound   = apperror.NewNotFoundError("user not found")
)
