package service

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ServerService struct {
	repository Repository
}

var _ UserClientService = (*ServerService)(nil)
var _ SyncService = (*ServerService)(nil)

func NewServerService(repository Repository) *ServerService {
	return &ServerService{
		repository: repository,
	}
}

func (s *ServerService) Register(ctx context.Context, login, password string) (domain.User, error) {
	ePassword, err := auth.EncryptPassword(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("auth.EncryptPassword: %w", err)
	}
	newUser := domain.NewUser(login, ePassword)
	user, err := s.repository.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("repository.CreateUser: %w", err)
	}
	return user, nil
}

func (s *ServerService) Login(ctx context.Context, login, password string) (string, error) {
	us, err := s.repository.FindUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("repository.FindUserByLogin: %w", err)
	}
	err = auth.ComparePasswords(us.HashedPassword, password)
	if err != nil {
		return "", err
	}
	token, err := auth.GenerateToken(us.ID)
	if err != nil {
		return "", fmt.Errorf("auth.GenerateToken: %w", err)
	}
	return token, nil
}

func (s *ServerService) RegisterClient(ctx context.Context) (domain.Client, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return domain.Client{}, fmt.Errorf("auth.GetUserID: %w", err)
	}
	id := uuid.New()
	client := domain.Client{
		ID:     id.String(),
		UserID: userID,
	}
	client, err = s.repository.CreateClient(ctx, client)
	if err != nil {
		return domain.Client{}, fmt.Errorf("repository.CreateClient: %w", err)
	}
	return client, nil
}

func (s *ServerService) CheckClient(ctx context.Context, id string) (domain.Client, error) {
	client, err := s.repository.FindClientByID(ctx, id)
	if err != nil {
		return domain.Client{}, fmt.Errorf("repository.FindClientByID: %w", err)
	}
	return client, nil
}

func (s *ServerService) SyncCredentials(ctx context.Context,
	sync *domain.CredSync) ([]domain.Credentials, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	creds := sync.Credentials
	log.Debug(fmt.Sprintf("credentials length for sync = %d", len(creds)))

	var modifiedAfter []domain.Credentials
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(creds); i++ {
			creds[i].UserID = userID
			errSave := s.repository.SaveCredentials(ctx, creds[i])
			if errSave != nil {
				log.Error(fmt.Sprintf("credentials with id = %s",
					creds[i].ID), zap.Error(errSave))
				return errSave
			}
		}
		modifiedAfter, err = s.repository.FindCredentialsModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return modifiedAfter, nil
}
func (s *ServerService) SyncCard(ctx context.Context,
	sync *domain.CardSync) ([]domain.Card, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	cards := sync.Cards
	log.Debug(fmt.Sprintf("cards length for sync = %d", len(cards)))

	var modifiedAfter []domain.Card
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(cards); i++ {
			cards[i].UserID = userID
			errSave := s.repository.SaveCard(ctx, cards[i])
			if errSave != nil {
				log.Error(fmt.Sprintf("card with id = %s",
					cards[i].ID), zap.Error(errSave))
				return errSave
			}
		}
		modifiedAfter, err = s.repository.FindCardsModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return modifiedAfter, nil
}
func (s *ServerService) SyncText(ctx context.Context, sync *domain.TextSync) ([]*domain.Text, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	texts := sync.Texts
	log.Debug(fmt.Sprintf("texts length for sync = %d", len(texts)))

	var textsAfter []*domain.Text
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(texts); i++ {
			texts[i].UserID = userID
			errSave := s.repository.SaveText(ctx, texts[i])
			if errSave != nil {
				log.Error(fmt.Sprintf("text with id = %s",
					texts[i].ID), zap.Error(errSave))
				return errSave
			}
		}
		textsAfter, err = s.repository.FindActiveTextsModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}

		textsDeleted, err := s.repository.FindDeletedTextsModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}
		textsAfter = append(textsAfter, textsDeleted...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return textsAfter, nil
}
func (s *ServerService) SyncBinary(ctx context.Context,
	sync *domain.BinarySync) ([]*domain.Binary, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	binaries := sync.Binaries
	log.Debug(fmt.Sprintf("binaries length for sync = %d", len(binaries)))

	var binaryAfter []*domain.Binary
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(binaries); i++ {
			binaries[i].UserID = userID
			errSave := s.repository.SaveBinary(ctx, binaries[i])
			if errSave != nil {
				log.Error(fmt.Sprintf("binary with id = %s",
					binaries[i].ID), zap.Error(errSave))
				return errSave
			}
		}
		binaryAfter, err = s.repository.FindActiveBinariesModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}

		binariesDeleted, err := s.repository.FindDeletedBinariesModifiedAfter(ctx, userID,
			sync.LastSyncTms)
		if err != nil {
			return err
		}
		binaryAfter = append(binaryAfter, binariesDeleted...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return binaryAfter, nil
}

func (s *ServerService) UpdateClientLastSyncTms(ctx context.Context, client domain.Client) error {
	err := s.repository.UpdateClientLastSyncTmsByID(ctx, client.ID, client.SyncTms)
	if err != nil {
		return fmt.Errorf("repository.UpdateClientLastSyncTmsByID: %w", err)
	}
	return nil
}
