package main

import (
	"github.com/denis-oreshkevich/gophkeeper/internal/server/server"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"go.uber.org/zap"
)

func main() {
	if err := server.Run(); err != nil {
		logger.Log.Fatal("main error", zap.Error(err))
	}
	logger.Log.Info("server exited properly")
}
