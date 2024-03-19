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

func DoCard(ctx context.Context, conf *config.Config, clientService *service.ClientService,
	dealer *crypto.Dealer, user model.User) error {
	switch conf.Action {
	case config.ActionGet:
		byID, err := clientService.FindCardByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.FindCardByID: %w", err)
		}
		num, err := dealer.Decrypt(byID.Num)
		if err != nil {
			return fmt.Errorf("dealer.Decrypt(byID.Num): %w", err)
		}
		cvc, err := dealer.Decrypt(byID.CVC)
		if err != nil {
			return fmt.Errorf("dealer.Decrypt(byID.CVC): %w", err)
		}
		name, err := dealer.Decrypt(byID.HolderName)
		if err != nil {
			return fmt.Errorf("dealer.Decrypt(byID.HolderName): %w", err)
		}
		logger.Log.Info(fmt.Sprintf("num: %s, cvc: %s, name %s",
			num, cvc, name))
	case config.ActionSave:
		id := uuid.New()
		num, err := dealer.Encrypt(conf.CardNum)
		if err != nil {
			return fmt.Errorf("dealer.Encrypt(conf.CardNum): %w", err)
		}
		cvc, err := dealer.Encrypt(conf.CardCVC)
		if err != nil {
			return fmt.Errorf("dealer.Encrypt(conf.CardCVC): %w", err)
		}
		name, err := dealer.Encrypt(conf.CardHolderName)
		if err != nil {
			return fmt.Errorf("dealer.Encrypt(conf.CardHolderName): %w", err)
		}
		card := model.Card{
			ID:          id.String(),
			Num:         num,
			CVC:         cvc,
			HolderName:  name,
			New:         conf.IsNew,
			UserID:      user.ID,
			Status:      model.StatusActive,
			ModifiedTms: time.Now().UTC(),
		}
		err = clientService.SaveCard(ctx, card)
		if err != nil {
			return fmt.Errorf("clientService.SaveCard: %w", err)
		}
		logger.Log.Info(fmt.Sprintf("saved card id = %s", id.String()))
	case config.ActionDelete:
		err := clientService.DeleteCardByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.DeleteCardByID: %w", err)
		}
		logger.Log.Info("Success")
	default:
		return fmt.Errorf("action value %s is unsupported", conf.Action)
	}
	return nil
}
