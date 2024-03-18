package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"time"
)

func (r *Repository) SaveCard(ctx context.Context, card *model.Card) error {
	if err := r.cardRepo.save(ctx, card); err != nil {
		return fmt.Errorf("cardRepo.save: %w", err)
	}
	return nil
}
func (r *Repository) FindCardByID(ctx context.Context, id string) (*model.Card, error) {
	byID, err := r.cardRepo.findByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cardRepo.findByID: %w", err)
	}
	return byID, nil
}
func (r *Repository) FindCardsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Card, error) {
	mod, err := r.cardRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("cardRepo.findActiveModifiedAfter: %w", err)
	}
	return mod, nil
}

func (r *Repository) DeleteCardByID(ctx context.Context, id string) error {
	if err := r.cardRepo.deleteByID(ctx, id); err != nil {
		return fmt.Errorf("cardRepo.deleteByID: %w", err)
	}
	return nil
}
