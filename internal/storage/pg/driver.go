package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/thxhix/chat/internal/config"
)

// PostgresStorage wraps a *sql.DB driver used by repository implementations.
type PostgresStorage struct {
	Driver *sql.DB
}

// OpenConnection opens a connection to Postgres using config.PostgresQL and verifies
// the connection by calling PingContext with the provided ctx.
//
// The returned *PostgresStorage should be closed by calling Driver.Close()
// (for example via a close function returned by the storage wiring).
func OpenConnection(ctx context.Context, cfg *config.Config) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &PostgresStorage{
		Driver: db,
	}, nil
}

type Migrator struct {
	db  *sql.DB
	src string
}

func NewMigrator(db *sql.DB, sourceURL string) *Migrator {
	return &Migrator{db: db, src: sourceURL}
}

func (m *Migrator) Up() error {
	driver, err := postgres.WithInstance(m.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate driver: %w", err)
	}
	mig, err := migrate.NewWithDatabaseInstance(m.src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}
