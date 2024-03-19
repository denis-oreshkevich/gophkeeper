package repo

import (
	"context"
	"errors"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"time"
)

var ErrUserAlreadyExist = errors.New("user already exist")

var ErrItemNotFound = errors.New("item not found")

// ClientRepository interface to access data
type ClientRepository interface {
	CreateUser(ctx context.Context, usr model.User) (model.User, error)
	FindUserByLogin(ctx context.Context, login string) (model.User, error)

	CreateClient(ctx context.Context, client model.Client) (model.Client, error)
	UpdateClientLastSyncTmsByID(ctx context.Context, id string, syncTms time.Time) error
	FindClientByID(ctx context.Context, id string) (model.Client, error)

	SaveCredentials(ctx context.Context, cred *model.Credentials) error
	FindCredentialsByID(ctx context.Context, id string) (*model.Credentials, error)
	FindCredentialsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Credentials, error)
	DeleteCredentialsByID(ctx context.Context, id string) error

	SaveText(ctx context.Context, txt *model.Text) error
	FindTextByID(ctx context.Context, id string) (*model.Text, error)
	FindActiveTextsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Text, error)
	FindDeletedTextsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Text, error)
	DeleteTextByID(ctx context.Context, id string) error

	SaveBinary(ctx context.Context, bin *model.Binary) error
	FindBinaryByID(ctx context.Context, id string) (*model.Binary, error)
	FindActiveBinariesModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Binary, error)
	FindDeletedBinariesModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Binary, error)
	DeleteBinaryByID(ctx context.Context, id string) error

	SaveCard(ctx context.Context, card *model.Card) error
	FindCardByID(ctx context.Context, id string) (*model.Card, error)
	FindCardsModifiedAfter(ctx context.Context, userID string,
		tms time.Time) ([]*model.Card, error)
	DeleteCardByID(ctx context.Context, id string) error

	InTransaction(ctx context.Context, transact func(context.Context) error) error
}
