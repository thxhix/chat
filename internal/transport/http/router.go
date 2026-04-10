package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/chat/internal/transport/http/middleware"
)

type RouteRegistrar interface {
	GetAPIRoutes(r chi.Router)
}

func NewRouter(handlers ...RouteRegistrar) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.TimeoutMiddleware)

	router.Route("/api", func(api chi.Router) {
		for _, h := range handlers {
			h.GetAPIRoutes(api)
		}
	})

	return router
}
