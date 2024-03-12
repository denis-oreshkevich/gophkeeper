package service

import (
	"context"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/adapter/api/client"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/auth"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type ClientService struct {
	baseRepo   Repository
	remoteRepo client.RESTRepository
}

var _ UserClientService = (*ClientService)(nil)

var _ SyncService = (*ClientService)(nil)

func NewClientService(baseRepo Repository, remoteRepo client.RESTRepository) *ClientService {
	return &ClientService{
		baseRepo:   baseRepo,
		remoteRepo: remoteRepo,
	}
}

func (s *ClientService) Register(ctx context.Context, login, password string) (domain.User, error) {
	ePassword, err := auth.EncryptPassword(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("auth.EncryptPassword: %w", err)
	}
	newUser := domain.NewUser(login, ePassword)
	newUser, err = s.remoteRepo.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("remoteRepo.CreateUser: %w", err)
	}

	newUser, err = s.baseRepo.CreateUser(ctx, newUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("baseRepo.CreateUser: %w", err)
	}

	return newUser, nil
}

func (s *ClientService) Login(ctx context.Context, login, password string) (string, error) {
	us, err := s.baseRepo.FindUserByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("baseRepo.FindUserByLogin: %w", err)
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

func (s *ClientService) RegisterClient(ctx context.Context) (domain.Client, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return domain.Client{}, fmt.Errorf("auth.GetUserID: %w", err)
	}
	id := uuid.New()
	client := domain.Client{
		ID:      id.String(),
		UserID:  userID,
		SyncTms: time.Now(),
	}
	client, err = s.remoteRepo.CreateClient(ctx, client)
	if err != nil {
		return domain.Client{}, fmt.Errorf("remoteRepo.CreateClient: %w", err)
	}

	client, err = s.baseRepo.CreateClient(ctx, client)
	return client, err
}

func (s *ClientService) CheckClient(ctx context.Context, id string) (domain.Client, error) {
	client, err := s.baseRepo.FindClientByID(ctx, id)
	if err != nil {
		return domain.Client{}, fmt.Errorf("repository.FindClientByID: %w", err)
	}
	return client, nil
}

func (s *ClientService) SyncCredentials(ctx context.Context,
	sync *domain.CredSync) ([]domain.Credentials, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
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
		err := s.baseRepo.SaveCredentials(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveCredentials: %w", err)
		}
	}
	return credentials, nil
}
func (s *ClientService) SyncCard(ctx context.Context,
	sync *domain.CardSync) ([]domain.Card, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
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
		err := s.baseRepo.SaveCard(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("baseRepo.SaveCard: %w", err)
		}
	}
	return cards, nil
}
func (s *ClientService) SyncText(ctx context.Context,
	sync *domain.TextSync) ([]*domain.Text, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
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
	sync *domain.BinarySync) ([]*domain.Binary, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth.GetUserID: %w", err)
	}
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

func (s *ClientService) UpdateClientLastSyncTms(ctx context.Context, client domain.Client) error {
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
