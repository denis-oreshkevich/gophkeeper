package service

import (
	"context"
	"errors"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"time"
)

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrItemNotFound = errors.New("item not found")

// Repository interface to access data
type Repository interface {
	CreateUser(ctx context.Context, usr domain.User) (domain.User, error)
	FindUserByLogin(ctx context.Context, login string) (domain.User, error)

	CreateClient(ctx context.Context, client domain.Client) (domain.Client, error)
	UpdateClientLastSyncTmsByID(ctx context.Context, id string, syncTms time.Time) error
	FindClientByID(ctx context.Context, id string) (domain.Client, error)
	FindClientsByUserID(ctx context.Context, userID string) ([]domain.Client, error)

	SaveCredentials(ctx context.Context, cred domain.Credentials) error
	FindCredentialsByID(ctx context.Context, id string) (domain.Credentials, error)
	FindCredentialsByUserID(ctx context.Context, userID string) ([]domain.Credentials, error)
	FindCredentialsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]domain.Credentials, error)
	DeleteCredentialsByID(ctx context.Context, id string) error

	SaveText(ctx context.Context, txt *domain.Text) error
	FindTextByID(ctx context.Context, id string) (*domain.Text, error)
	FindTextsByUserID(ctx context.Context, userID string) ([]*domain.Text, error)
	FindActiveTextsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*domain.Text, error)
	FindDeletedTextsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*domain.Text, error)
	DeleteTextByID(ctx context.Context, id string) error

	SaveBinary(ctx context.Context, bin *domain.Binary) error
	FindBinaryByID(ctx context.Context, id string) (*domain.Binary, error)
	FindBinariesByUserID(ctx context.Context, userID string) ([]*domain.Binary, error)
	FindActiveBinariesModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*domain.Binary, error)
	FindDeletedBinariesModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*domain.Binary, error)
	DeleteBinaryByID(ctx context.Context, id string) error

	SaveCard(ctx context.Context, card domain.Card) error
	FindCardByID(ctx context.Context, id string) (domain.Card, error)
	FindCardsByUserID(ctx context.Context, userID string) ([]domain.Card, error)
	FindCardsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]domain.Card, error)
	DeleteCardByID(ctx context.Context, id string) error

	InTransaction(ctx context.Context, transact func(context.Context) error) error
}
