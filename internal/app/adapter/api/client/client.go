package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/config"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type RESTRepository interface {
	CreateUser(ctx context.Context, usr domain.User) (domain.User, error)
	CreateClient(ctx context.Context, client domain.Client) (domain.Client, error)
	UpdateClientLastSyncTms(ctx context.Context, id string, syncTms time.Time) error

	SyncCredentials(ctx context.Context, sync *domain.CredSync) ([]domain.Credentials, error)
	SyncCard(ctx context.Context, sync *domain.CardSync) ([]domain.Card, error)
	SyncText(ctx context.Context, sync *domain.TextSync) ([]*domain.Text, error)
	SyncBinary(ctx context.Context, sync *domain.BinarySync) ([]*domain.Binary, error)
}

type RESTRepositoryImpl struct {
	client *resty.Client
	conf   *config.Config
}

var _ RESTRepository = (*RESTRepositoryImpl)(nil)

func NewRESTRepositoryImpl(client *resty.Client, conf *config.Config) *RESTRepositoryImpl {
	return &RESTRepositoryImpl{
		client: client,
		conf:   conf,
	}
}

func (r RESTRepositoryImpl) CreateUser(ctx context.Context, usr domain.User) (domain.User, error) {
	marshal, err := json.Marshal(usr)
	if err != nil {
		return domain.User{}, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/register`)
	if err != nil {
		return domain.User{}, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusOK {
		return domain.User{}, fmt.Errorf("response status code = %d", status)
	}

	body := response.Body()
	err = json.Unmarshal(body, &usr)
	if err != nil {
		return domain.User{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return usr, nil
}

func (r RESTRepositoryImpl) CreateClient(ctx context.Context,
	client domain.Client) (domain.Client, error) {
	marshal, err := json.Marshal(client)
	if err != nil {
		return domain.Client{}, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/client`)
	if err != nil {
		return domain.Client{}, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusCreated {
		return domain.Client{}, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	err = json.Unmarshal(body, &client)
	if err != nil {
		return domain.Client{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return client, nil
}

func (r RESTRepositoryImpl) UpdateClientLastSyncTms(ctx context.Context, id string, syncTms time.Time) error {
	client := domain.Client{
		ID:      id,
		SyncTms: syncTms,
	}

	marshal, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Put(r.conf.ServerAddress + `/api/user/client`)
	if err != nil {
		return fmt.Errorf("client.R().Put: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("response status code = %d", status)
	}
	return nil
}

func (r RESTRepositoryImpl) SyncCredentials(ctx context.Context,
	sync *domain.CredSync) ([]domain.Credentials, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/credentials/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dc []domain.Credentials
	err = json.Unmarshal(body, &dc)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dc, nil
}

func (r RESTRepositoryImpl) SyncCard(ctx context.Context,
	sync *domain.CardSync) ([]domain.Card, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/cards/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dc []domain.Card
	err = json.Unmarshal(body, &dc)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dc, nil
}

func (r RESTRepositoryImpl) SyncText(ctx context.Context,
	sync *domain.TextSync) ([]*domain.Text, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/texts/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dt []*domain.Text
	err = json.Unmarshal(body, &dt)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dt, nil
}

func (r RESTRepositoryImpl) SyncBinary(ctx context.Context,
	sync *domain.BinarySync) ([]*domain.Binary, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.conf.ServerAddress + `/api/user/binaries/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var db []*domain.Binary
	err = json.Unmarshal(body, &db)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return db, nil
}
