package fs

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"time"
)

func (r *Repository) SaveBinary(ctx context.Context, bin *model.Binary) error {
	err := r.binaryRepo.save(ctx, bin)
	if err != nil {
		return fmt.Errorf("binaryRepo.save: %w", err)
	}
	return nil
}

func (r *Repository) FindBinaryByID(ctx context.Context, id string) (*model.Binary, error) {
	byID, err := r.binaryRepo.findByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("binaryRepo.findByID: %w", err)
	}
	return byID, nil
}

func (r *Repository) FindBinariesByUserID(ctx context.Context,
	userID string) ([]*model.Binary, error) {
	binary, err := r.binaryRepo.findByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("binaryRepo.findByUserID: %w", err)
	}
	return binary, nil
}

func (r *Repository) FindActiveBinariesModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Binary, error) {
	binary, err := r.binaryRepo.findActiveModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("binaryRepo.findActiveModifiedAfter: %w", err)
	}
	return binary, nil
}

func (r *Repository) FindDeletedBinariesModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Binary, error) {
	binary, err := r.binaryRepo.findDeletedModifiedAfter(ctx, userID, tms)
	if err != nil {
		return nil, fmt.Errorf("binaryRepo.findDeletedModifiedAfter: %w", err)
	}
	return binary, nil
}

func (r *Repository) DeleteBinaryByID(ctx context.Context, id string) error {
	if err := r.binaryRepo.deleteByID(ctx, id); err != nil {
		return fmt.Errorf("binaryRepo.deleteByID: %w", err)
	}
	return nil
}
