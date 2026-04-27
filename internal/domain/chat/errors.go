package chat

import "github.com/thxhix/chat/internal/apperror"

var (
	ErrNotFound  = apperror.NewNotFoundError("Chat not found")
	ErrForbidden = apperror.NewForbiddenError("You don't have access to chat")

	ErrEmptyMembersProvided   = apperror.NewBadRequestError("You must provide at least one member")
	ErrTooManyMembersProvided = apperror.NewBadRequestError("Too many members provided")

	ErrBadDirectTypeProvided       = apperror.NewBadRequestError("Bad direct type provided")
	ErrPrivateDirectTooManyMembers = apperror.NewBadRequestError("Too many members provided for private chat")
)
