package storage

import (
	"context"
	"github.com/thxhix/chat/internal/config"
	"github.com/thxhix/chat/internal/domain/token"
	"github.com/thxhix/chat/internal/domain/user"
	"github.com/thxhix/chat/internal/logger"
	"github.com/thxhix/chat/internal/storage/pg"
	tokenpg "github.com/thxhix/chat/internal/storage/pg/token"
	userpg "github.com/thxhix/chat/internal/storage/pg/user"

	"go.uber.org/zap"
)

// Storage groups repository instances used by the application.
type Storage struct {
	User  user.IUserRepository
	Token token.ITokenRepository
}

// NewStorage creates repository instances, runs migrations and returns a cleanup
// function. It expects cfg.PostgresQL to be set to a valid DSN. If cfg.PostgresQL
// is empty ErrNoPostgresConnection is returned.
//
// The returned close function should be called to close DB connections when the
// application shuts down.
//
// Example:
//
//	storage, closeFn, err := storage.NewStorage(ctx, cfg, logger)
//	if err != nil { return err }
//	defer closeFn()
func NewStorage(ctx context.Context, cfg *config.Config, logger logger.ILogger) (*Storage, func(), error) {
	if cfg.DatabaseURI == "" {
		return nil, nil, ErrNoPostgresConnection
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.DatabaseInitTimeout)
	defer cancel()

	logger.Info("Trying to connect to postgresql", zap.String("DSN", cfg.DatabaseURI))
	db, err := pg.OpenConnection(ctx, cfg)
	if err != nil {
		return nil, nil, err
	}

	logger.Info("Trying to migrate", zap.String("migrations_path", cfg.MigrationsPath))
	migrator := pg.NewMigrator(db.Driver, cfg.MigrationsPath)
	if err := migrator.Up(); err != nil {
		return nil, nil, err
	}

	closeFn := func() { _ = db.Driver.Close() }

	userRepository := userpg.NewRepository(db.Driver)
	tokenRepository := tokenpg.NewRepository(db.Driver)

	return &Storage{
		User:  userRepository,
		Token: tokenRepository,
	}, closeFn, nil
}
