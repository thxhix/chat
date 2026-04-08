package token

import (
	"crypto/subtle"
	"time"
)

// ValidateToken checks whether a refresh token is valid for a given user.
//
// It performs several checks:
//  1. The token's user ID must match the provided userId.
//  2. The token must not have been replaced (ReplacedBy is nil).
//  3. The token must not be expired (based on now).
//  4. The token hash must match the provided incomingHash using constant-time comparison.
//
// Returns an error if any check fails.
func ValidateToken(now time.Time, userId int64, incomingHash string, rec *RefreshTokenRecord) error {
	if rec.UserID != userId {
		return ErrInvalidCredentials
	}
	if rec.ReplacedBy != nil {
		return ErrInvalidCredentials
	}
	if now.After(rec.ExpiresAt) {
		return ErrInvalidCredentials
	}
	if subtle.ConstantTimeCompare([]byte(rec.TokenHash), []byte(incomingHash)) != 1 {
		return ErrInvalidCredentials
	}
	return nil
}

// ValidateRefreshToken checks that a refresh token string is not empty.
//
// Returns ErrEmptyRefreshCredentials if the token is empty.
func ValidateRefreshToken(token string) error {
	if token == "" {
		return ErrInvalidCredentials
	}
	return nil
}
