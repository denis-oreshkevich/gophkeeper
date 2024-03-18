package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"go.uber.org/zap"
)

type ServerService struct {
	repository repo.ServerRepository
}

func NewServerService(repository repo.ServerRepository) ServerService {
	return ServerService{
		repository: repository,
	}
}

func (s *ServerService) Register(ctx context.Context, login, password string) (model.User, error) {
	ePassword, err := auth.EncryptPassword(password)
	if err != nil {
		return model.User{}, fmt.Errorf("auth.EncryptPassword: %w", err)
	}
	newUser := model.NewUser(login, ePassword)
	user, err := s.repository.CreateUser(ctx, newUser)
	if err != nil {
		return model.User{}, fmt.Errorf("repository.CreateUser: %w", err)
	}
	user.HashedPassword = ""
	return user, nil
}

func (s *ServerService) Login(ctx context.Context, login, password string) (model.User, error) {
	us, err := s.repository.FindUserByLogin(ctx, login)
	if err != nil {
		return model.User{}, fmt.Errorf("repository.FindUserByLogin: %w", err)
	}
	err = auth.ComparePasswords(us.HashedPassword, password)
	if err != nil {
		return model.User{}, err
	}
	us.HashedPassword = ""
	return us, nil
}

func (s *ServerService) RegisterClient(ctx context.Context, client model.Client) (model.Client, error) {
	client, err := s.repository.CreateClient(ctx, client)
	if err != nil {
		return model.Client{}, fmt.Errorf("repository.CreateClient: %w", err)
	}
	return client, nil
}

func (s *ServerService) CheckClient(ctx context.Context, id string) (model.Client, error) {
	client, err := s.repository.FindClientByID(ctx, id)
	if err != nil {
		return model.Client{}, fmt.Errorf("repository.FindClientByID: %w", err)
	}
	return client, nil
}

func (s *ServerService) SaveCredentials(ctx context.Context, cred model.Credentials) error {
	err := s.repository.SaveCredentials(ctx, cred)
	if err != nil {
		return fmt.Errorf("repository.SaveCredentials: %w", err)
	}
	return nil
}

func (s *ServerService) FindCredentialsByID(ctx context.Context, id string) (model.Credentials, error) {
	cred, err := s.repository.FindCredentialsByID(ctx, id)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("repository.FindCredentialsByID: %w", err)
	}
	return cred, nil

}

func (s *ServerService) FindCredentialsByUserID(ctx context.Context) ([]model.Credentials, error) {
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

func (s *ServerService) DeleteCredentialsByID(ctx context.Context, id string) error {
	err := s.repository.DeleteCredentialsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteCredentialsByID: %w", err)
	}
	return nil
}

func (s *ServerService) SaveText(ctx context.Context, txt *model.Text) error {
	err := s.repository.SaveText(ctx, txt)
	if err != nil {
		return fmt.Errorf("repository.SaveText: %w", err)
	}
	return nil
}

func (s *ServerService) FindTextByID(ctx context.Context, id string) (*model.Text, error) {
	text, err := s.repository.FindTextByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.FindTextByID: %w", err)
	}
	return text, nil
}

func (s *ServerService) FindTextsByUserID(ctx context.Context) ([]*model.Text, error) {
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

func (s *ServerService) DeleteTextByID(ctx context.Context, id string) error {
	err := s.repository.DeleteTextByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteTextByID: %w", err)
	}
	return nil
}

func (s *ServerService) SaveBinary(ctx context.Context, bin *model.Binary) error {
	err := s.repository.SaveBinary(ctx, bin)
	if err != nil {
		return fmt.Errorf("repository.SaveBinary: %w", err)
	}
	return nil
}

func (s *ServerService) FindBinaryByID(ctx context.Context, id string) (*model.Binary, error) {
	b, err := s.repository.FindBinaryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repository.FindBinaryByID: %w", err)
	}
	return b, nil
}

func (s *ServerService) FindBinariesByUserID(ctx context.Context) ([]*model.Binary, error) {
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

func (s *ServerService) DeleteBinaryByID(ctx context.Context, id string) error {
	err := s.repository.DeleteBinaryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteBinaryByID: %w", err)
	}
	return nil
}

func (s *ServerService) SaveCard(ctx context.Context, card model.Card) error {
	err := s.repository.SaveCard(ctx, card)
	if err != nil {
		return fmt.Errorf("repository.SaveCard: %w", err)
	}
	return nil
}

func (s *ServerService) FindCardByID(ctx context.Context, id string) (model.Card, error) {
	card, err := s.repository.FindCardByID(ctx, id)
	if err != nil {
		return model.Card{}, fmt.Errorf("repository.FindCardByID: %w", err)
	}
	return card, nil
}

func (s *ServerService) FindCardsByUserID(ctx context.Context) ([]model.Card, error) {
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

func (s *ServerService) DeleteCardByID(ctx context.Context, id string) error {
	err := s.repository.DeleteCardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("repository.DeleteCardByID: %w", err)
	}
	return nil
}

func (s *ServerService) SyncCredentials(ctx context.Context,
	sync *model.CredSync) ([]model.Credentials, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	creds := sync.Credentials
	log.Debug(fmt.Sprintf("credentials length for sync = %d", len(creds)))

	modifiedAfter, err := s.repository.FindCredentialsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindCredentialsModifiedAfter: %w", err)
	}

	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(creds); i++ {
			credentials := creds[i]
			credentials.UserID = userID
			saved, err := s.repository.FindCredentialsByID(ctx, credentials.ID)
			var checkNeeded = true
			if err != nil {
				if !errors.Is(err, repo.ErrItemNotFound) {
					return fmt.Errorf("repository.FindCredentialsByID: %w", err)
				}
				checkNeeded = false
			}
			if checkNeeded && saved.ModifiedTms.After(credentials.ModifiedTms) {
				log.Debug(fmt.Sprintf("credentials with id = %s "+
					"did not saved^ because newer version was saved", credentials.ID))
				continue
			}
			errSave := s.repository.SaveCredentials(ctx, *credentials)
			if errSave != nil {
				log.Error(fmt.Sprintf("credentials with id = %s",
					credentials.ID), zap.Error(errSave))
				return errSave
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return modifiedAfter, nil
}
func (s *ServerService) SyncCard(ctx context.Context,
	sync *model.CardSync) ([]model.Card, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	cards := sync.Cards
	log.Debug(fmt.Sprintf("cards length for sync = %d", len(cards)))

	modifiedAfter, err := s.repository.FindCardsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindCardsModifiedAfter: %w", err)
	}
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(cards); i++ {
			card := cards[i]
			card.UserID = userID
			saved, err := s.repository.FindCardByID(ctx, card.ID)
			var checkNeeded = true
			if err != nil {
				if !errors.Is(err, repo.ErrItemNotFound) {
					return fmt.Errorf("repository.FindCardByID: %w", err)
				}
				checkNeeded = false
			}
			if checkNeeded && saved.ModifiedTms.After(card.ModifiedTms) {
				log.Debug(fmt.Sprintf("card with id = %s "+
					"did not saved^ because newer version was saved", card.ID))
				continue
			}
			errSave := s.repository.SaveCard(ctx, *card)
			if errSave != nil {
				log.Error(fmt.Sprintf("card with id = %s",
					card.ID), zap.Error(errSave))
				return errSave
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return modifiedAfter, nil
}
func (s *ServerService) SyncText(ctx context.Context, sync *model.TextSync) ([]*model.Text, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	texts := sync.Texts
	log.Debug(fmt.Sprintf("texts length for sync = %d", len(texts)))

	textsAfter, err := s.repository.FindActiveTextsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindActiveTextsModifiedAfter: %w", err)
	}

	textsDeleted, err := s.repository.FindDeletedTextsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindDeletedTextsModifiedAfter: %w", err)
	}
	textsAfter = append(textsAfter, textsDeleted...)
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(texts); i++ {
			text := texts[i]
			text.UserID = userID
			saved, err := s.repository.FindTextByID(ctx, text.ID)
			var checkNeeded = true
			if err != nil {
				if !errors.Is(err, repo.ErrItemNotFound) {
					return fmt.Errorf("repository.FindTextByID: %w", err)
				}
				checkNeeded = false
			}
			if checkNeeded && saved.ModifiedTms.After(text.ModifiedTms) {
				log.Debug(fmt.Sprintf("card with id = %s "+
					"did not saved^ because newer version was saved", text.ID))
				continue
			}
			errSave := s.repository.SaveText(ctx, text)
			if errSave != nil {
				log.Error(fmt.Sprintf("text with id = %s",
					text.ID), zap.Error(errSave))
				return errSave
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return textsAfter, nil
}
func (s *ServerService) SyncBinary(ctx context.Context,
	sync *model.BinarySync) ([]*model.Binary, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
	log := logger.Log.With(zap.String("userID", userID))

	binaries := sync.Binaries
	log.Debug(fmt.Sprintf("binaries length for sync = %d", len(binaries)))

	binaryAfter, err := s.repository.FindActiveBinariesModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindActiveBinariesModifiedAfter: %w", err)
	}

	binariesDeleted, err := s.repository.FindDeletedBinariesModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("repository.FindDeletedBinariesModifiedAfter: %w", err)
	}
	binaryAfter = append(binaryAfter, binariesDeleted...)
	err = s.repository.InTransaction(ctx, func(ctx context.Context) error {
		for i := 0; i < len(binaries); i++ {
			binary := binaries[i]
			binary.UserID = userID
			saved, err := s.repository.FindBinaryByID(ctx, binary.ID)
			var checkNeeded = true
			if err != nil {
				if !errors.Is(err, repo.ErrItemNotFound) {
					return fmt.Errorf("repository.FindBinaryByID: %w", err)
				}
				checkNeeded = false
			}
			if checkNeeded && saved.ModifiedTms.After(binary.ModifiedTms) {
				log.Debug(fmt.Sprintf("card with id = %s "+
					"did not saved^ because newer version was saved", binary.ID))
				continue
			}
			errSave := s.repository.SaveBinary(ctx, binary)
			if errSave != nil {
				log.Error(fmt.Sprintf("binary with id = %s",
					binary.ID), zap.Error(errSave))
				return errSave
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repository.InTransaction: %w", err)
	}
	return binaryAfter, nil
}

func (s *ServerService) UpdateClientLastSyncTms(ctx context.Context, client model.Client) error {
	err := s.repository.UpdateClientLastSyncTmsByID(ctx, client.ID, client.SyncTms)
	if err != nil {
		return fmt.Errorf("repository.UpdateClientLastSyncTmsByID: %w", err)
	}
	return nil
}
