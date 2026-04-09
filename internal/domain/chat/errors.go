package chat

import "github.com/thxhix/chat/internal/apperror"

var (
	ErrNotFound  = apperror.NewNotFoundError("Chat not found")
	ErrForbidden = apperror.NewForbiddenError("You don't have access to chat")
)
