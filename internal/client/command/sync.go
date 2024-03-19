package command

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/service"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"time"
)

func DoSync(ctx context.Context, findClient model.Client,
	clientService *service.ClientService, userID string) error {
	now := time.Now().UTC()
	binarySync := model.BinarySync{LastSyncTms: findClient.SyncTms}
	_, err := clientService.SyncBinary(ctx, &binarySync, userID)
	if err != nil {
		return fmt.Errorf("clientService.SyncBinary: %w", err)
	}

	cardSync := model.CardSync{LastSyncTms: findClient.SyncTms}
	_, err = clientService.SyncCard(ctx, &cardSync, userID)
	if err != nil {
		return fmt.Errorf("clientService.SyncCard: %w", err)
	}

	credSync := model.CredSync{LastSyncTms: findClient.SyncTms}
	_, err = clientService.SyncCredentials(ctx, &credSync, userID)
	if err != nil {
		return fmt.Errorf("clientService.SyncCredentials: %w", err)
	}

	textSync := model.TextSync{LastSyncTms: findClient.SyncTms}
	_, err = clientService.SyncText(ctx, &textSync, userID)
	if err != nil {
		return fmt.Errorf("clientService.SyncText: %w", err)
	}

	findClient.SyncTms = now
	err = clientService.UpdateClientLastSyncTms(ctx, findClient)
	if err != nil {
		return fmt.Errorf("clientService.UpdateClientLastSyncTms: %w", err)
	}
	return nil
}
