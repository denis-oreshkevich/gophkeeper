package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"time"
)

func (r *Repository) SaveCredentials(ctx context.Context, cred domain.Credentials) error {
	if err := r.credentialsRepo.save(ctx, cred); err != nil {
		return fmt.Errorf("credentialsRepo.save: %w", err)
	}
	return nil
}
func (r *Repository) FindCredentialsByID(ctx context.Context,
	id string) (domain.Credentials, error) {
	byID, err := r.credentialsRepo.findByID(ctx, id)
	if err != nil {
		return domain.Credentials{}, fmt.Errorf("credentialsRepo.findByID: %w", err)
	}
	return byID, nil

}
func (r *Repository) FindCredentialsByUserID(ctx context.Context,
	userID string) ([]domain.Credentials, error) {
	byID, err := r.credentialsRepo.findByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("credentialsRepo.findByUserID: %w", err)
	}
	return byID, nil
}
func (r *Repository) FindCredentialsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]domain.Credentials, error) {
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
