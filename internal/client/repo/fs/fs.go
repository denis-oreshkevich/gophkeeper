package fs

import (
	"context"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
)

type Repository struct {
	userRepo        *UserRepository
	clientRepo      *ClientRepository
	binaryRepo      *BaseRepository[*model.Binary]
	cardRepo        *BaseRepository[*model.Card]
	credentialsRepo *BaseRepository[*model.Credentials]
	textRepo        *BaseRepository[*model.Text]
}

func NewRepository(userRepo *UserRepository, clientRepo *ClientRepository,
	binaryRepo *BaseRepository[*model.Binary],
	cardRepo *BaseRepository[*model.Card],
	credentialsRepo *BaseRepository[*model.Credentials],
	textRepo *BaseRepository[*model.Text]) *Repository {
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
