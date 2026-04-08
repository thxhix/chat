package token

import (
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	JTI        uuid.UUID
	UserID     uuid.UUID
	TokenHash  string
	IssuedAt   time.Time
	ExpiresAt  time.Time
	ReplacedBy *uuid.UUID
	CreatedAt  time.Time
}

type RefreshTokenRecord struct {
	JTI        uuid.UUID
	UserID     int64
	TokenHash  string
	IssuedAt   time.Time
	ExpiresAt  time.Time
	ReplacedBy *uuid.UUID
}
