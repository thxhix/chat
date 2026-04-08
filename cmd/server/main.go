package main

import (
	"github.com/thxhix/chat/internal/app"
	"github.com/thxhix/chat/internal/config"
	"github.com/thxhix/chat/internal/logger"
)

func main() {
	log, err := logger.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Close() }()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	err = app.RunServer(log, cfg)
	if err != nil {
		panic(err)
	}
}
