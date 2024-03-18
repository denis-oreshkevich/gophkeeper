package command

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/crypto"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/google/uuid"
	"time"
)

func DoCredentials(ctx context.Context, conf *config.Config,
	clientService *service.ClientService, dealer *crypto.Dealer,
	user model.User) error {
	switch conf.Action {
	case config.ActionGet:
		byID, err := clientService.FindCredentialsByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.FindCredentialsByID: %w", err)
		}
		login, err := dealer.Decrypt(byID.Login)
		if err != nil {
			return fmt.Errorf("dealer.Decrypt(byID.Login): %w", err)
		}
		password, err := dealer.Decrypt(byID.Password)
		if err != nil {
			return fmt.Errorf("dealer.Decrypt(byID.Password): %w", err)
		}

		logger.Log.Info(fmt.Sprintf("Login: %s, password: %s",
			login, password))
	case config.ActionSave:
		id := uuid.New()
		crLogin, err := dealer.Encrypt(conf.CredentialsLogin)
		if err != nil {
			return fmt.Errorf("dealer.Encrypt(conf.CredentialsLogin): %w", err)
		}

		crPwd, err := dealer.Encrypt(conf.CredentialsPassword)
		if err != nil {
			return fmt.Errorf("dealer.Encrypt(conf.CredentialsPassword): %w", err)
		}
		cred := model.Credentials{
			ID:          id.String(),
			Login:       crLogin,
			Password:    crPwd,
			New:         conf.IsNew,
			UserID:      user.ID,
			Status:      model.StatusActive,
			ModifiedTms: time.Now().UTC(),
		}
		err = clientService.SaveCredentials(ctx, cred)
		if err != nil {
			return fmt.Errorf("clientService.SaveCredentials: %w", err)
		}
		logger.Log.Info(fmt.Sprintf("saved credentials id = %s", id.String()))
	case config.ActionDelete:
		err := clientService.DeleteCredentialsByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.DeleteCredentialsByID: %w", err)
		}
		logger.Log.Info("Success")
	default:
		return fmt.Errorf("action value %s is unsupported", conf.Action)
	}
	return nil
}
