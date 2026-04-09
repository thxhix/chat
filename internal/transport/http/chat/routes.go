package chat

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/chat/internal/transport/http/middleware"
)

const CHAT_ID_KEY = "chat_id"

func (h *Handler) GetAPIRoutes(r chi.Router) {
	r.Route("/chats", func(r chi.Router) {
		r.Use(middleware.Authorize(h.logger, h.jwtManager))
		r.Route("/{"+CHAT_ID_KEY+"}", func(r chi.Router) {
			r.Get("/messages", h.GetMessages)
			r.Post("/messages", h.SendMessage)
		})
	})
}
