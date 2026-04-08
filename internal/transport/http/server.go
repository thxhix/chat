package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thxhix/chat/internal/config"
	"github.com/thxhix/chat/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// server implements the Server interface.
type Server struct {
	server *http.Server
	router *chi.Mux
	cfg    *config.Config
	logger logger.ILogger
}

// NewServer creates a new Server instance with the provided router, config, and logger.
// The returned Server is ready to start with the Start() method.
func NewServer(router *chi.Mux, cfg *config.Config, logger logger.ILogger) *Server {
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	return &Server{
		server: srv,
		router: router,
		cfg:    cfg,
		logger: logger,
	}
}

// Start runs the HTTP server in a separate goroutine and waits for either:
// 1. A server error from ListenAndServe (returns that error), or
// 2. A termination signal (SIGINT, SIGTERM, SIGQUIT) to gracefully shutdown the server.
// It returns an error if the shutdown fails.
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("HTTP server startup", zap.String("address", s.cfg.Address))

	errCh := make(chan error, 1)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(sigCh)

	select {
	case err := <-errCh:
		return err

	case sig := <-sigCh:
		s.logger.Info("HTTP server graceful shutdown by signal", zap.String("signal", sig.String()))
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server error", zap.Error(err))
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	}
}
