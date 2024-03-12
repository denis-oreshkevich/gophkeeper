package fs

import (
	"context"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
)

type Repository struct {
	userRepo        *UserRepository
	clientRepo      *ClientRepository
	binaryRepo      *BaseRepository[*domain.Binary]
	cardRepo        *BaseRepository[*domain.Card]
	credentialsRepo *BaseRepository[*domain.Credentials]
	textRepo        *BaseRepository[*domain.Text]
}

func NewRepository(userRepo *UserRepository, clientRepo *ClientRepository,
	binaryRepo *BaseRepository[*domain.Binary],
	cardRepo *BaseRepository[*domain.Card],
	credentialsRepo *BaseRepository[*domain.Credentials],
	textRepo *BaseRepository[*domain.Text]) *Repository {
	return &Repository{
		userRepo:        userRepo,
		clientRepo:      clientRepo,
		binaryRepo:      binaryRepo,
		cardRepo:        cardRepo,
		credentialsRepo: credentialsRepo,
		textRepo:        textRepo,
	}
}

func (r *Repository) InTransaction(ctx context.Context,
	transact func(context.Context) error) error {
	err := transact(ctx)
	logger.Log.Debug("transaction not supported for file repository")
	return err
}
