package command

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

var buildVersion = "N/A"

var buildDate = "N/A"

var buildCommit = "N/A"

func Run() error {
	ctx := context.Background()
	err := logger.Initialize(zapcore.DebugLevel.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "logger initialize", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	logger.Log.Info(fmt.Sprintf("Build version: %s\n", buildVersion))
	logger.Log.Info(fmt.Sprintf("Build date: %s\n", buildDate))
	logger.Log.Info(fmt.Sprintf("Build commit: %s\n", buildCommit))

	conf, err := config.Parse()
	if err != nil {
		logger.Log.Fatal("parse config", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	err = Do(ctx, conf)
	return err
}
