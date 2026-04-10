package handlers

import (
	"github.com/thxhix/chat/internal/transport/http/auth"
	"github.com/thxhix/chat/internal/transport/http/chat"
	"github.com/thxhix/chat/internal/transport/http/message"
)

type Handlers struct {
	Auth    *auth.Handler
	Chat    *chat.Handler
	Message *message.Handler
}
