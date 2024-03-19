package command

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo/fs"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo/rest"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/crypto"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"os"
	"strings"
)

func Do(ctx context.Context, conf *config.Config) error {
	err := os.MkdirAll(conf.WorkingDir, 0755)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("os.MkdirAll: %w", err)
		}
	}

	userRepo := fs.NewUserRepository(conf.WorkingDir +
		string(os.PathSeparator) + "user.json")

	clientRepo := fs.NewClientRepository(conf.WorkingDir+
		string(os.PathSeparator)+"client.json",
		"client-tm-*")

	binaryRepo := fs.NewBaseRepository(conf.WorkingDir+
		string(os.PathSeparator)+"binary.json",
		"binary-tm-*", &model.Binary{})

	cardRepo := fs.NewBaseRepository(conf.WorkingDir+
		string(os.PathSeparator)+"card.json",
		"card-tm-*", &model.Card{})

	credentialsRepo := fs.NewBaseRepository(conf.WorkingDir+
		string(os.PathSeparator)+"credentials.json",
		"credentials-tm-*", &model.Credentials{})

	textRepo := fs.NewBaseRepository(conf.WorkingDir+
		string(os.PathSeparator)+"text.json",
		"text-tm-*", &model.Text{})

	repository := fs.NewRepository(userRepo, clientRepo, binaryRepo,
		cardRepo, credentialsRepo, textRepo)

	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/key.pem")
	if err != nil {
		return fmt.Errorf("tls.LoadX509KeyPair: %w", err)
	}
	r := resty.New().SetCertificates(cert).SetTLSClientConfig(&tls.Config{
		// на маке не доверяет
		InsecureSkipVerify: true,
	}).SetBaseURL("https://" + conf.ServerAddress)

	restRepo := rest.NewRESTRepositoryImpl(r)

	clientService := service.NewClientService(repository, restRepo)

	user, err := login(ctx, conf, repository, clientService)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return fmt.Errorf("auth.GenerateToken: %w", err)
	}

	r.SetAuthToken(token)

	findClient, err := clientRepo.FindClient(ctx)
	if err != nil {
		if !errors.Is(err, repo.ErrItemNotFound) {
			return fmt.Errorf("clientRepo.FindClient: %w", err)
		}
		findClient, err = clientService.RegisterClient(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("clientService.RegisterClient: %w", err)
		}
	}

	if !conf.IsSync && conf.Action == "" {
		return errors.New("action is empty and")
	}
	k := user.ID + user.Login

	dealer, err := crypto.NewDealer(k)
	if err != nil {
		return fmt.Errorf("crypto.NewDealer %w", err)
	}

	logger.Log.Debug(fmt.Sprintf("id is %s", conf.ID))

	if conf.IsFileFlagsParsed {
		err := DoFile(ctx, conf, clientService, user)
		if err != nil {
			return fmt.Errorf("DoFile: %w", err)
		}
	} else if conf.IsTextFlagsParsed {
		err := DoText(ctx, conf, clientService, user)
		if err != nil {
			return fmt.Errorf("DoText: %w", err)
		}
	} else if conf.IsCardFlagsParsed {
		err := DoCard(ctx, conf, clientService, dealer, user)
		if err != nil {
			return fmt.Errorf("DoCard: %w", err)
		}

	} else if conf.IsCredentialsFlagsParsed {
		err := DoCredentials(ctx, conf, clientService, dealer, user)
		if err != nil {
			return fmt.Errorf("DoCredentials: %w", err)
		}
	} else if conf.IsSync {
		err = DoSync(ctx, findClient, clientService, user.ID)
		if err != nil {
			return fmt.Errorf("DoSync: %w", err)
		}
	} else {
		logger.Log.Error("nothing to do", zap.Error(err))
	}
	return nil
}

func login(ctx context.Context, conf *config.Config, repository *fs.Repository,
	clientService *service.ClientService) (model.User, error) {
	user, err := repository.FindUser(ctx)

	if err != nil {
		if !errors.Is(err, repo.ErrItemNotFound) {
			return model.User{}, fmt.Errorf("repository.FindUser: %w", err)
		}
		l := conf.UserLogin
		password := conf.UserPassword
		if strings.TrimSpace(l) == "" || strings.TrimSpace(password) == "" {
			return model.User{}, errors.New("invalid credentials")
		}
		user, err = clientService.Login(ctx, l, password)
		if err != nil {
			if !errors.Is(err, repo.ErrItemNotFound) {
				return model.User{}, fmt.Errorf("clientService.Login: %w", err)
			}
			user, err = clientService.RegisterUser(ctx, l, password)
			if err != nil {
				return model.User{}, fmt.Errorf("clientService.RegisterUser: %w", err)
			}
		}

	}
	return user, nil
}
