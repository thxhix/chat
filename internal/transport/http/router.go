package http

import (
	"github.com/go-chi/chi/v5"
)

type RouteRegistrar interface {
	GetAPIRoutes(r chi.Router)
}

func NewRouter(handlers ...RouteRegistrar) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(api chi.Router) {
		for _, h := range handlers {
			h.GetAPIRoutes(api)
		}
	})

	return router
}
