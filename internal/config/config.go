package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/thxhix/chat/internal/security/jwt"
	"github.com/thxhix/chat/internal/storage/pg/core/tx_manager"
	"os"
	"time"
)

type Config struct {
	Address             string        `env:"ADDRESS" envDefault:"localhost:8080"`
	DatabaseURI         string        `env:"DATABASE_URI"`
	DatabaseInitTimeout time.Duration `env:"DB_INIT_TIMEOUT" envDefault:"15s"`
	MigrationsPath      string        `env:"MIGRATIONS_PATH" envDefault:"file://migrations"`

	UUIDSalt     string           `env:"UUID_SALT" envDefault:"chat"`
	TXManagerKey tx_manager.TXKey `env:"TX_MANAGER_KEY" envDefault:"chat"`

	// Embedded JWT configuration.
	JWT jwt.JWTConfig
}

// Priority: flags > env > defaults
func (c *Config) parseFlags(args []string) error {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)

	addr := fs.String("addr", "", "Server address")
	dbURI := fs.String("db", "", "DB connection string")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if *addr != "" {
		c.Address = *addr
	}
	if *dbURI != "" {
		c.DatabaseURI = *dbURI
	}
	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.parseFlags(os.Args[1:])
	if err != nil {
		return nil, err
	}

	if cfg.DatabaseURI == "" {
		return nil, ErrDatabaseURIMissing
	}

	return cfg, nil
}
