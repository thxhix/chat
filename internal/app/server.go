package app

import (
	"context"
	"github.com/thxhix/chat/internal/config"
	authdomain "github.com/thxhix/chat/internal/domain/auth"
	chatdomain "github.com/thxhix/chat/internal/domain/chat"
	messagedomain "github.com/thxhix/chat/internal/domain/message"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/security"
	"github.com/thxhix/chat/internal/storage"
	"github.com/thxhix/chat/internal/transport/http"
	authhttp "github.com/thxhix/chat/internal/transport/http/auth"
	chathttp "github.com/thxhix/chat/internal/transport/http/chat"
	"go.uber.org/zap"
)

func RunServer(logger logger.ILogger, cfg *config.Config) error {
	ctx := context.Background()

	store, closeFn, err := storage.NewStorage(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to create repo storage", zap.Error(err))
		return err
	}
	defer closeFn()

	sec := security.NewSecurity(cfg)

	// Services
	as := authdomain.NewService(store.User, store.Token, sec.Password, sec.JWT, sec.Crypt)
	cs := chatdomain.NewService(store.Chat, store.ChatMember)
	ms := messagedomain.NewMessageService(cs, store.Message)

	// Handlers
	ah := authhttp.NewHandler(logger, as)
	ch := chathttp.NewHandler(logger, sec.JWT, ms)

	r := http.NewRouter(ah, ch)

	s := http.NewServer(r, cfg, logger)
	err = s.Start(ctx)
	if err != nil {
		logger.Error("Failed to start http server", zap.Error(err))
		return err
	}

	return err
}
