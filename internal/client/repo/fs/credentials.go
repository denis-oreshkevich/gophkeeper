package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"time"
)

func (r *Repository) SaveCredentials(ctx context.Context, cred *model.Credentials) error {
	if err := r.credentialsRepo.save(ctx, cred); err != nil {
		return fmt.Errorf("credentialsRepo.save: %w", err)
	}
	return nil
}
func (r *Repository) FindCredentialsByID(ctx context.Context,
	id string) (*model.Credentials, error) {
	byID, err := r.credentialsRepo.findByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("credentialsRepo.findByID: %w", err)
	}
	return byID, nil

}

func (r *Repository) FindCredentialsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Credentials, error) {
	mod, err := r.credentialsRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("credentialsRepo.findActiveModifiedAfter: %w", err)
	}
	return mod, nil
}
func (r *Repository) DeleteCredentialsByID(ctx context.Context, id string) error {
	if err := r.credentialsRepo.deleteByID(ctx, id); err != nil {
		return fmt.Errorf("credentialsRepo.deleteByID: %w", err)
	}
	return nil
}
