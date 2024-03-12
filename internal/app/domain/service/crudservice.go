package service

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
)

type CRUDServiceImpl struct {
	repository Repository
}

var _ CRUDService = (*CRUDServiceImpl)(nil)

func NewCRUDServiceImpl(repository Repository) *CRUDServiceImpl {
	return &CRUDServiceImpl{
		repository: repository,
	}
}

func (s *CRUDServiceImpl) SaveCredentials(ctx context.Context, cred domain.Credentials) error {
	err := s.repository.SaveCredentials(ctx, cred)
	if err != nil {
		return fmt.Errorf("repository.SaveCredentials: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) FindCredentialsByID(ctx context.Context, id string) (domain.Credentials, error) {
	cred, err := s.repository.FindCredentialsByID(ctx, id)
	if err != nil {
		return domain.Credentials{}, fmt.Errorf("repository.FindCredentialsByID: %w", err)
	}
	return cred, nil

}

func (s *CRUDServiceImpl) FindCredentialsByUserID(ctx context.Context) ([]domain.Credentials, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	creds, err := s.repository.FindCredentialsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.FindCredentialsByUserID: %w", err)
	}
	return creds, nil
}

func (s *CRUDServiceImpl) DeleteCredentialsByID(ctx context.Context, id string) error {
	err := s.repository.DeleteCredentialsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteCredentialsByID: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) SaveText(ctx context.Context, txt *domain.Text) error {
	err := s.repository.SaveText(ctx, txt)
	if err != nil {
		return fmt.Errorf("repository.SaveText: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) FindTextByID(ctx context.Context, id string) (*domain.Text, error) {
	text, err := s.repository.FindTextByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.FindTextByID: %w", err)
	}
	return text, nil
}

func (s *CRUDServiceImpl) FindTextsByUserID(ctx context.Context) ([]*domain.Text, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	texts, err := s.repository.FindTextsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.FindTextsByUserID: %w", err)
	}
	return texts, nil
}

func (s *CRUDServiceImpl) DeleteTextByID(ctx context.Context, id string) error {
	err := s.repository.DeleteTextByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteTextByID: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) SaveBinary(ctx context.Context, bin *domain.Binary) error {
	err := s.repository.SaveBinary(ctx, bin)
	if err != nil {
		return fmt.Errorf("repository.SaveBinary: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) FindBinaryByID(ctx context.Context, id string) (*domain.Binary, error) {
	b, err := s.repository.FindBinaryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.FindBinaryByID: %w", err)
	}
	return b, nil
}

func (s *CRUDServiceImpl) FindBinariesByUserID(ctx context.Context) ([]*domain.Binary, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	binaries, err := s.repository.FindBinariesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.FindBinariesByUserID: %w", err)
	}
	return binaries, nil
}

func (s *CRUDServiceImpl) DeleteBinaryByID(ctx context.Context, id string) error {
	err := s.repository.DeleteBinaryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteBinaryByID: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) SaveCard(ctx context.Context, card domain.Card) error {
	err := s.repository.SaveCard(ctx, card)
	if err != nil {
		return fmt.Errorf("repository.SaveCard: %w", err)
	}
	return nil
}

func (s *CRUDServiceImpl) FindCardByID(ctx context.Context, id string) (domain.Card, error) {
	card, err := s.repository.FindCardByID(ctx, id)
	if err != nil {
		return domain.Card{}, fmt.Errorf("repository.FindCardByID: %w", err)
	}
	return card, nil
}

func (s *CRUDServiceImpl) FindCardsByUserID(ctx context.Context) ([]domain.Card, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	cards, err := s.repository.FindCardsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.FindCardsByUserID: %w", err)
	}
	return cards, nil
}

func (s *CRUDServiceImpl) DeleteCardByID(ctx context.Context, id string) error {
	err := s.repository.DeleteCardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteCardByID: %w", err)
	}
	return nil
}
