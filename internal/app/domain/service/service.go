package service

import (
	"context"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
)

type UserClientService interface {
	Register(ctx context.Context, login, password string) (domain.User, error)
	Login(ctx context.Context, login, password string) (string, error)

	RegisterClient(ctx context.Context) (domain.Client, error)
	CheckClient(ctx context.Context, id string) (domain.Client, error)
}

type CRUDService interface {
	SaveCredentials(ctx context.Context, cred domain.Credentials) error
	FindCredentialsByID(ctx context.Context, id string) (domain.Credentials, error)
	FindCredentialsByUserID(ctx context.Context) ([]domain.Credentials, error)
	DeleteCredentialsByID(ctx context.Context, id string) error

	SaveText(ctx context.Context, txt *domain.Text) error
	FindTextByID(ctx context.Context, id string) (*domain.Text, error)
	FindTextsByUserID(ctx context.Context) ([]*domain.Text, error)
	DeleteTextByID(ctx context.Context, id string) error

	SaveBinary(ctx context.Context, bin *domain.Binary) error
	FindBinaryByID(ctx context.Context, id string) (*domain.Binary, error)
	FindBinariesByUserID(ctx context.Context) ([]*domain.Binary, error)
	DeleteBinaryByID(ctx context.Context, id string) error

	SaveCard(ctx context.Context, card domain.Card) error
	FindCardByID(ctx context.Context, id string) (domain.Card, error)
	FindCardsByUserID(ctx context.Context) ([]domain.Card, error)
	DeleteCardByID(ctx context.Context, id string) error
}

type SyncService interface {
	SyncCredentials(ctx context.Context, sync *domain.CredSync) ([]domain.Credentials, error)
	SyncCard(ctx context.Context, sync *domain.CardSync) ([]domain.Card, error)
	SyncText(ctx context.Context, sync *domain.TextSync) ([]*domain.Text, error)
	SyncBinary(ctx context.Context, sync *domain.BinarySync) ([]*domain.Binary, error)

	UpdateClientLastSyncTms(ctx context.Context, client domain.Client) error
}
