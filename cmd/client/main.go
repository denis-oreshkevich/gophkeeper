package main

import (
	"github.com/denis-oreshkevich/gophkeeper/internal/client/command"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"go.uber.org/zap"
)

func main() {
	if err := command.Run(); err != nil {
		logger.Log.Fatal("main error", zap.Error(err))
	}
}
