package app

import (
	"context"
	"github.com/thxhix/chat/internal/config"
	authdomain "github.com/thxhix/chat/internal/domain/auth"
	chatdomain "github.com/thxhix/chat/internal/domain/chat"
	messagedomain "github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/security"
	"github.com/thxhix/chat/internal/security/uuid"
	"github.com/thxhix/chat/internal/storage"
	"github.com/thxhix/chat/internal/transport/http"
	authhttp "github.com/thxhix/chat/internal/transport/http/auth"
	chathttp "github.com/thxhix/chat/internal/transport/http/chat"
	"github.com/thxhix/chat/internal/transport/http/core/handlers"
	"github.com/thxhix/chat/internal/transport/http/core/router"
	msghttp "github.com/thxhix/chat/internal/transport/http/message"
	"go.uber.org/zap"
)

func RunServer(logger logger.ILogger, cfg *config.Config) error {
	ctx := context.Background()

	sec := security.NewSecurity(cfg)
	uuidManager := uuid.NewUUIDManager(cfg.UUIDSalt)

	store, driver, txManager, err := storage.NewStorage(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to create repo storage", zap.Error(err))
		return err
	}
	defer func() { _ = driver.Close() }()

	// Services
	as := authdomain.NewService(store.User, store.Token, sec.Password, sec.JWT, sec.Crypt)
	cs := chatdomain.NewService(store.Chat, store.ChatMember, txManager, uuidManager)
	ms := messagedomain.NewMessageService(cs, store.Message, uuidManager)

	// Handlers
	h := &handlers.Handlers{
		Auth:    authhttp.NewHandler(logger, as),
		Message: msghttp.NewHandler(logger, ms),
		Chat:    chathttp.NewHandler(logger, ms, cs),
	}

	r := router.NewRouter(logger, sec.JWT, h)

	s := http.NewServer(r, cfg, logger)
	err = s.Start(ctx)
	if err != nil {
		logger.Error("Failed to start http server", zap.Error(err))
		return err
	}

	return err
}
