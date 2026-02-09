package main

import (
	"github.com/gersastas/wallet-service/internal/config"
	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	cfg := config.New(logger)
	server := httpserver.New(cfg.GetHTTPBindAddr())

	if err := server.Run(); err != nil {
		logger.Panic("HTTP server failed", zap.Error(err))
	}
}
