package command

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/google/uuid"
	"time"
)

func DoText(ctx context.Context, conf *config.Config, clientService *service.ClientService,
	user model.User) error {
	switch conf.Action {
	case config.ActionGet:
		byID, err := clientService.FindTextByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.FindTextByID: %w", err)
		}
		logger.Log.Info(byID.Txt)
	case config.ActionSave:
		id := uuid.New()
		txt := model.Text{
			ID:          id.String(),
			Txt:         conf.Text,
			New:         conf.IsNew,
			UserID:      user.ID,
			Status:      model.StatusActive,
			ModifiedTms: time.Now().UTC(),
		}
		err := clientService.SaveText(ctx, &txt)
		if err != nil {
			return fmt.Errorf("clientService.SaveText: %w", err)
		}
		logger.Log.Info(fmt.Sprintf("saved text id = %s", id.String()))
	case config.ActionDelete:
		err := clientService.DeleteTextByID(ctx, conf.ID)
		if err != nil {
			return fmt.Errorf("clientService.DeleteTextByID: %w", err)
		}
		logger.Log.Info("Success")
	default:
		return fmt.Errorf("action value %s is unsupported", conf.Action)
	}
	return nil
}
