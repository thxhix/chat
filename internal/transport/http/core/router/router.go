package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/transport/http/core"
	"github.com/thxhix/chat/internal/transport/http/core/handlers"
	"github.com/thxhix/chat/internal/transport/http/middleware"
)

func NewRouter(logger logger.ILogger, jwtManager jwt.IJWTManager, h *handlers.Handlers) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.NewRecoverer(logger))
	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.TimeoutMiddleware)

	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authorize(logger, jwtManager))
			r.Route("/chats", func(r chi.Router) {
				r.Post("/", h.Chat.CreateChat)

				r.Route("/{"+core.ChatIdPrefix+"}", func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Get("/messages", h.Message.GetMessages)
						r.Post("/messages", h.Message.SendMessage)
					})
				})
			})
		})
	})

	return router
}

//func walkRoutes(r chi.Router) {
//	fmt.Println("\n--- Registered Routes ---")
//
//	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
//		// Убираем лишние слеши в конце для красоты
//		route = strings.Replace(route, "/*", "/", -1)
//
//		// Выводим метод (зеленым или просто выровненным) и сам путь
//		fmt.Printf("%-7s %s\n", method, route)
//		return nil
//	}
//
//	if err := chi.Walk(r, walkFunc); err != nil {
//		fmt.Printf("Logging error: %s\n", err.Error())
//	}
//
//	fmt.Println("-------------------------\n")
//}
