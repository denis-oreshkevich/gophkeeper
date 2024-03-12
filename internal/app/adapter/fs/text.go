package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"time"
)

func (r *Repository) SaveText(ctx context.Context, txt *domain.Text) error {
	err := r.textRepo.save(ctx, txt)
	if err != nil {
		return fmt.Errorf("textRepo.save: %w", err)
	}
	return nil
}

func (r *Repository) FindTextByID(ctx context.Context, id string) (*domain.Text, error) {
	byID, err := r.textRepo.findByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("textRepo.findByID: %w", err)
	}
	return byID, nil
}

func (r *Repository) FindTextsByUserID(ctx context.Context,
	userID string) ([]*domain.Text, error) {
	byID, err := r.textRepo.findByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("textRepo.findByUserID: %w", err)
	}
	return byID, nil
}
func (r *Repository) FindActiveTextsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*domain.Text, error) {
	after, err := r.textRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("textRepo.findActiveModifiedAfter: %w", err)
	}
	return after, nil
}

func (r *Repository) FindDeletedTextsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*domain.Text, error) {
	after, err := r.textRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("textRepo.findActiveModifiedAfter: %w", err)
	}
	return after, nil
}

func (r *Repository) DeleteTextByID(ctx context.Context, id string) error {
	if err := r.textRepo.deleteByID(ctx, id); err != nil {
		return fmt.Errorf("textRepo.deleteByID: %w", err)
	}
	return nil
}
