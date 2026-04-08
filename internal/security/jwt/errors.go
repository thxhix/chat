package jwt

import (
	"errors"
	"github.com/thxhix/chat/internal/apperror"
)

var (
	ErrSecretTooShort          = errors.New("JWT secret is too short")
	ErrUnexpectedSigningMethod = errors.New("JWT unexpected signing method")

	ErrInvalidAccessToken  = errors.New("JWT access token expired")
	ErrInvalidRefreshToken = apperror.NewUnauthorizedError("invalid refresh token")
)
