package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo/rest"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/logger"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type ClientService struct {
	baseRepo   repo.ClientRepository
	remoteRepo rest.RESTRepository
}

func NewClientService(baseRepo repo.ClientRepository,
	remoteRepo rest.RESTRepository) *ClientService {
	return &ClientService{
		baseRepo:   baseRepo,
		remoteRepo: remoteRepo,
	}
}

func (s *ClientService) RegisterUser(ctx context.Context, login,
	password string) (model.User, error) {
	ausr := model.AuthUser{
		Login:    login,
		Password: password,
	}
	usr, err := s.remoteRepo.CreateUser(ctx, ausr)
	if err != nil {
		return model.User{}, fmt.Errorf("remoteRepo.CreateUser: %w", err)
	}

	ePassword, err := auth.EncryptPassword(password)
	newUser := model.NewUser(login, ePassword)
	newUser.ID = usr.ID
	if err != nil {
		return model.User{}, fmt.Errorf("auth.EncryptPassword: %w", err)
	}
	newUser, err = s.baseRepo.CreateUser(ctx, newUser)
	if err != nil {
		return model.User{}, fmt.Errorf("baseRepo.CreateUser: %w", err)
	}

	return newUser, nil
}

func (s *ClientService) Login(ctx context.Context, login, password string) (model.User, error) {
	us, err := s.baseRepo.FindUserByLogin(ctx, login)
	if err != nil {
		if !errors.Is(err, repo.ErrItemNotFound) {
			return model.User{}, fmt.Errorf("baseRepo.FindUserByLogin: %w", err)
		}
		authUser := model.AuthUser{
			Login:    login,
			Password: password,
		}
		user, err := s.remoteRepo.Login(ctx, authUser)
		if err != nil {
			return model.User{}, fmt.Errorf("remoteRepo.Login: %w", err)
		}
		us.ID = user.ID
		us.Login = login
		ePassword, err := auth.EncryptPassword(password)
		if err != nil {
			return model.User{}, fmt.Errorf("auth.EncryptPassword: %w", err)
		}
		us.HashedPassword = ePassword
		us, err = s.baseRepo.CreateUser(ctx, us)
		if err != nil {
			return model.User{}, fmt.Errorf("baseRepo.CreateUser: %w", err)
		}
	}
	err = auth.ComparePasswords(us.HashedPassword, password)
	if err != nil {
		return model.User{}, err
	}
	return us, nil
}

func (s *ClientService) RegisterClient(ctx context.Context, userID string) (model.Client, error) {
	id := uuid.New()
	client := model.Client{
		ID:      id.String(),
		UserID:  userID,
		SyncTms: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	client, err := s.remoteRepo.CreateClient(ctx, client)
	if err != nil {
		return model.Client{}, fmt.Errorf("remoteRepo.CreateClient: %w", err)
	}

	client, err = s.baseRepo.CreateClient(ctx, client)
	return client, err
}

func (s *ClientService) CheckClient(ctx context.Context, id string) (model.Client, error) {
	client, err := s.baseRepo.FindClientByID(ctx, id)
	if err != nil {
		return model.Client{}, fmt.Errorf("baseRepo.FindClientByID: %w", err)
	}
	return client, nil
}

func (s *ClientService) SaveCredentials(ctx context.Context, cred model.Credentials) error {
	err := s.baseRepo.SaveCredentials(ctx, &cred)
	if err != nil {
		return fmt.Errorf("baseRepo.SaveCredentials: %w", err)
	}
	return nil
}

func (s *ClientService) FindCredentialsByID(ctx context.Context, id string) (model.Credentials, error) {
	cred, err := s.baseRepo.FindCredentialsByID(ctx, id)
	if err != nil {
		return model.Credentials{}, fmt.Errorf("baseRepo.FindCredentialsByID: %w", err)
	}
	return *cred, nil

}

func (s *ClientService) DeleteCredentialsByID(ctx context.Context, id string) error {
	err := s.baseRepo.DeleteCredentialsByID(ctx, id)
	if err != nil {
		return fmt.Errorf("baseRepo.DeleteCredentialsByID: %w", err)
	}
	return nil
}

func (s *ClientService) SaveText(ctx context.Context, txt *model.Text) error {
	err := s.baseRepo.SaveText(ctx, txt)
	if err != nil {
		return fmt.Errorf("baseRepo.SaveText: %w", err)
	}
	return nil
}

func (s *ClientService) FindTextByID(ctx context.Context, id string) (*model.Text, error) {
	text, err := s.baseRepo.FindTextByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindTextByID: %w", err)
	}
	return text, nil
}

func (s *ClientService) DeleteTextByID(ctx context.Context, id string) error {
	err := s.baseRepo.DeleteTextByID(ctx, id)
	if err != nil {
		return fmt.Errorf("baseRepo.DeleteTextByID: %w", err)
	}
	return nil
}

func (s *ClientService) SaveBinary(ctx context.Context, bin *model.Binary) error {
	err := s.baseRepo.SaveBinary(ctx, bin)
	if err != nil {
		return fmt.Errorf("baseRepo.SaveBinary: %w", err)
	}
	return nil
}

func (s *ClientService) FindBinaryByID(ctx context.Context, id string) (*model.Binary, error) {
	b, err := s.baseRepo.FindBinaryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindBinaryByID: %w", err)
	}
	return b, nil
}

func (s *ClientService) DeleteBinaryByID(ctx context.Context, id string) error {
	err := s.baseRepo.DeleteBinaryByID(ctx, id)
	if err != nil {
		return fmt.Errorf("baseRepo.DeleteBinaryByID: %w", err)
	}
	return nil
}

func (s *ClientService) SaveCard(ctx context.Context, card model.Card) error {
	err := s.baseRepo.SaveCard(ctx, &card)
	if err != nil {
		return fmt.Errorf("baseRepo.SaveCard: %w", err)
	}
	return nil
}

func (s *ClientService) FindCardByID(ctx context.Context, id string) (model.Card, error) {
	card, err := s.baseRepo.FindCardByID(ctx, id)
	if err != nil {
		return model.Card{}, fmt.Errorf("baseRepo.FindCardByID: %w", err)
	}
	return *card, nil
}

func (s *ClientService) DeleteCardByID(ctx context.Context, id string) error {
	err := s.baseRepo.DeleteCardByID(ctx, id)
	if err != nil {
		return fmt.Errorf("baseRepo.DeleteCardByID: %w", err)
	}
	return nil
}

func (s *ClientService) SyncCredentials(ctx context.Context,
	sync *model.CredSync, userID string) ([]model.Credentials, error) {
	log := logger.Log.With(zap.String("userID", userID))

	modifiedAfter, err := s.baseRepo.FindCredentialsModifiedAfter(ctx, userID,
		sync.LastSyncTms)

	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindCredentialsModifiedAfter: %w", err)
	}

	sync.Credentials = modifiedAfter

	log.Debug(fmt.Sprintf("credentials length for sync = %d", len(sync.Credentials)))

	credentials, err := s.remoteRepo.SyncCredentials(ctx, sync)
	if err != nil {
		return nil, fmt.Errorf("remoteRepo.SyncCredentials: %w", err)
	}

	for i := 0; i < len(credentials); i++ {
		c := credentials[i]
		c.New = true
		err := s.baseRepo.SaveCredentials(ctx, &c)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveCredentials: %w", err)
		}
	}
	return credentials, nil
}
func (s *ClientService) SyncCard(ctx context.Context,
	sync *model.CardSync, userID string) ([]model.Card, error) {
	log := logger.Log.With(zap.String("userID", userID))

	modifiedAfter, err := s.baseRepo.FindCardsModifiedAfter(ctx, userID,
		sync.LastSyncTms)

	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindCardsModifiedAfter: %w", err)
	}

	sync.Cards = modifiedAfter

	log.Debug(fmt.Sprintf("cards length for sync = %d", len(sync.Cards)))

	cards, err := s.remoteRepo.SyncCard(ctx, sync)
	if err != nil {
		return nil, fmt.Errorf("remoteRepo.SyncCard: %w", err)
	}
	for i := 0; i < len(cards); i++ {
		c := cards[i]
		c.New = true
		err := s.baseRepo.SaveCard(ctx, &c)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveCard: %w", err)
		}
	}
	return cards, nil
}
func (s *ClientService) SyncText(ctx context.Context,
	sync *model.TextSync, userID string) ([]*model.Text, error) {
	log := logger.Log.With(zap.String("userID", userID))

	textsAfter, err := s.baseRepo.FindActiveTextsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindActiveTextsModifiedAfter: %w", err)
	}

	textsDeleted, err := s.baseRepo.FindDeletedTextsModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindDeletedTextsModifiedAfter: %w", err)
	}
	textsAfter = append(textsAfter, textsDeleted...)

	sync.Texts = textsAfter

	log.Debug(fmt.Sprintf("texts length for sync = %d", len(sync.Texts)))

	texts, err := s.remoteRepo.SyncText(ctx, sync)
	if err != nil {
		return nil, fmt.Errorf("remoteRepo.SyncText: %w", err)
	}
	for i := 0; i < len(texts); i++ {
		t := texts[i]
		t.New = true
		err := s.baseRepo.SaveText(ctx, t)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveText: %w", err)
		}
	}
	return texts, nil
}
func (s *ClientService) SyncBinary(ctx context.Context,
	sync *model.BinarySync, userID string) ([]*model.Binary, error) {
	log := logger.Log.With(zap.String("userID", userID))

	binariesAfter, err := s.baseRepo.FindActiveBinariesModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindActiveBinariesModifiedAfter: %w", err)
	}

	binariesDeleted, err := s.baseRepo.FindDeletedBinariesModifiedAfter(ctx, userID,
		sync.LastSyncTms)
	if err != nil {
		return nil, fmt.Errorf("baseRepo.FindDeletedBinariesModifiedAfter: %w", err)
	}

	if err != nil {
		return nil, err
	}
	binariesAfter = append(binariesAfter, binariesDeleted...)

	sync.Binaries = binariesAfter

	log.Debug(fmt.Sprintf("binaries length for sync = %d", len(sync.Binaries)))

	binaries, err := s.remoteRepo.SyncBinary(ctx, sync)

	if err != nil {
		return nil, fmt.Errorf("remoteRepo.SyncBinary: %w", err)
	}
	for i := 0; i < len(binaries); i++ {
		b := binaries[i]
		b.New = true
		err := s.baseRepo.SaveBinary(ctx, b)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveBinary: %w", err)
		}
	}
	return binaries, nil
}

func (s *ClientService) UpdateClientLastSyncTms(ctx context.Context, client model.Client) error {
	err := s.remoteRepo.UpdateClientLastSyncTms(ctx, client.ID, client.SyncTms)
	if err != nil {
		return fmt.Errorf("remoteRepo.UpdateClientLastSyncTmsByID: %w", err)
	}

	err = s.baseRepo.UpdateClientLastSyncTmsByID(ctx, client.ID, client.SyncTms)
	if err != nil {
		return fmt.Errorf("baseRepo.UpdateClientLastSyncTmsByID: %w", err)
	}

	return nil
}
