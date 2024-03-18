package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/client/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type RESTRepository interface {
	Login(ctx context.Context, usr model.AuthUser) (model.User, error)
	CreateUser(ctx context.Context, usr model.AuthUser) (model.User, error)
	CreateClient(ctx context.Context, client model.Client) (model.Client, error)
	UpdateClientLastSyncTms(ctx context.Context, id string, syncTms time.Time) error

	SyncCredentials(ctx context.Context, sync *model.CredSync) ([]model.Credentials, error)
	SyncCard(ctx context.Context, sync *model.CardSync) ([]model.Card, error)
	SyncText(ctx context.Context, sync *model.TextSync) ([]*model.Text, error)
	SyncBinary(ctx context.Context, sync *model.BinarySync) ([]*model.Binary, error)
}

type RESTRepositoryImpl struct {
	client *resty.Client
}

var _ RESTRepository = (*RESTRepositoryImpl)(nil)

func NewRESTRepositoryImpl(client *resty.Client) *RESTRepositoryImpl {
	return &RESTRepositoryImpl{
		client: client,
	}
}

func (r RESTRepositoryImpl) Login(ctx context.Context, usr model.AuthUser) (model.User, error) {
	marshal, err := json.Marshal(usr)
	if err != nil {
		return model.User{}, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/login`)
	if err != nil {
		return model.User{}, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusOK {
		if status == http.StatusUnauthorized {
			return model.User{}, repo.ErrItemNotFound
		}
		return model.User{}, fmt.Errorf("response status code = %d", status)
	}

	body := response.Body()
	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return model.User{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return user, nil
}

func (r RESTRepositoryImpl) CreateUser(ctx context.Context, usr model.AuthUser) (model.User, error) {
	marshal, err := json.Marshal(usr)
	if err != nil {
		return model.User{}, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/register`)
	if err != nil {
		return model.User{}, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusOK {
		return model.User{}, fmt.Errorf("response status code = %d", status)
	}

	body := response.Body()
	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return model.User{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return user, nil
}

func (r RESTRepositoryImpl) CreateClient(ctx context.Context,
	client model.Client) (model.Client, error) {
	marshal, err := json.Marshal(client)
	if err != nil {
		return model.Client{}, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/client`)
	if err != nil {
		return model.Client{}, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusCreated {
		return model.Client{}, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	err = json.Unmarshal(body, &client)
	if err != nil {
		return model.Client{}, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return client, nil
}

func (r RESTRepositoryImpl) UpdateClientLastSyncTms(ctx context.Context, id string, syncTms time.Time) error {
	client := model.Client{
		ID:      id,
		SyncTms: syncTms,
	}

	marshal, err := json.Marshal(client)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Put(r.client.BaseURL + `/api/user/client`)
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
	sync *model.CredSync) ([]model.Credentials, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/credentials/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dc []model.Credentials
	err = json.Unmarshal(body, &dc)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dc, nil
}

func (r RESTRepositoryImpl) SyncCard(ctx context.Context,
	sync *model.CardSync) ([]model.Card, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/cards/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dc []model.Card
	err = json.Unmarshal(body, &dc)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dc, nil
}

func (r RESTRepositoryImpl) SyncText(ctx context.Context,
	sync *model.TextSync) ([]*model.Text, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/texts/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var dt []*model.Text
	err = json.Unmarshal(body, &dt)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return dt, nil
}

func (r RESTRepositoryImpl) SyncBinary(ctx context.Context,
	sync *model.BinarySync) ([]*model.Binary, error) {
	marshal, err := json.Marshal(sync)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	response, err := r.client.R().
		SetContext(ctx).SetBody(marshal).Post(r.client.BaseURL + `/api/user/binaries/sync`)
	if err != nil {
		return nil, fmt.Errorf("client.R().Post: %w", err)
	}
	status := response.StatusCode()
	if status != http.StatusAccepted {
		return nil, fmt.Errorf("response status code = %d", status)
	}
	body := response.Body()
	var db []*model.Binary
	err = json.Unmarshal(body, &db)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return db, nil
}
