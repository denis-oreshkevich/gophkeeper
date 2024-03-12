package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/adapter/api/client"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/adapter/fs"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var buildVersion = "N/A"

var buildDate = "N/A"

var buildCommit = "N/A"

func main() {
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

	userRepo := fs.NewUserRepository("user")
	clientRepo := fs.NewClientRepository("client")
	binaryRepo := fs.NewBaseRepository("binary",
		"binary-tm-*", &domain.Binary{})
	cardRepo := fs.NewBaseRepository("card",
		"card-tm-*", &domain.Card{})
	credentialsRepo := fs.NewBaseRepository("credentials",
		"credentials-tm-*", &domain.Credentials{})
	textRepo := fs.NewBaseRepository("text",
		"text-tm-*", &domain.Text{})

	repository := fs.NewRepository(userRepo, clientRepo, binaryRepo,
		cardRepo, credentialsRepo, textRepo)

	user, err := userRepo.FindUser(ctx)

	if err != nil {
		//login
	}
	token, err := auth.GenerateToken(user.ID)

	cert, err := tls.LoadX509KeyPair("./certs/cert.pem", "./certs/key.pem")
	rest := resty.New().SetCertificates(cert).SetAuthToken(token)

	restRepo := client.NewRESTRepositoryImpl(rest, conf)

	clientService := service.NewClientService(repository, restRepo)

	//login := conf.UserLogin
	//password := conf.UserPassword

	findClient, err := clientRepo.FindClient(ctx)
	if err != nil {

	}

	if conf.IsFileFlagsParsed {

	} else if conf.IsTextFlagsParsed {

	} else if conf.IsCardFlagsParsed {

	} else if conf.IsCredentialsFlagsParsed {

	} else if conf.Sync {
		now := time.Now()
		binarySync := domain.BinarySync{LastSyncTms: findClient.SyncTms}
		_, err := clientService.SyncBinary(ctx, &binarySync)
		if err != nil {
			return
		}

		cardSync := domain.CardSync{LastSyncTms: findClient.SyncTms}
		_, err = clientService.SyncCard(ctx, &cardSync)
		if err != nil {
			logger.Log.Error("clientService.SyncCard", zap.Error(err))
			return
		}

		credSync := domain.CredSync{LastSyncTms: findClient.SyncTms}
		_, err = clientService.SyncCredentials(ctx, &credSync)
		if err != nil {
			logger.Log.Error("clientService.SyncCredentials", zap.Error(err))
			return
		}

		textSync := domain.TextSync{LastSyncTms: findClient.SyncTms}
		_, err = clientService.SyncText(ctx, &textSync)
		if err != nil {
			logger.Log.Error("clientService.SyncText", zap.Error(err))
			return
		}

		findClient.SyncTms = now
		err = clientService.UpdateClientLastSyncTms(ctx, findClient)
		if err != nil {
			logger.Log.Error("clientService.UpdateClientLastSyncTms", zap.Error(err))
			return
		}
	} else {
		logger.Log.Error("nothing to do", zap.Error(err))
		stop()
	}

}
