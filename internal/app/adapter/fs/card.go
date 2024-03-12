package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"time"
)

func (r *Repository) SaveCard(ctx context.Context, card domain.Card) error {
	if err := r.cardRepo.save(ctx, card); err != nil {
		return fmt.Errorf("cardRepo.save")
	}
	return nil
}
func (r *Repository) FindCardByID(ctx context.Context, id string) (domain.Card, error) {
	byID, err := r.cardRepo.findByID(ctx, id)
	if err != nil {
		return domain.Card{}, fmt.Errorf("cardRepo.findByID: %w", err)
	}
	return byID, nil
}
func (r *Repository) FindCardsByUserID(ctx context.Context, userID string) ([]domain.Card, error) {
	byID, err := r.cardRepo.findByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cardRepo.findByUserID: %w", err)
	}
	return byID, nil
}
func (r *Repository) FindCardsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]domain.Card, error) {
	mod, err := r.cardRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("cardRepo.findActiveModifiedAfter: %w", err)
	}
	return mod, nil
}

func (r *Repository) DeleteCardByID(ctx context.Context, id string) error {
	if err := r.cardRepo.deleteByID(ctx, id); err != nil {
		return fmt.Errorf("cardRepo.deleteByID")
	}
	return nil
}
