package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r *Repository) CreateClient(ctx context.Context, client domain.Client) (domain.Client, error) {
	query := `insert into keeper.client(id, user_id) 
	values (@id, @user_id) returning client.sync_tms`
	args := pgx.NamedArgs{
		"id":      client.ID,
		"user_id": client.UserID,
	}
	row := r.db.QueryRow(ctx, query, args)
	err := row.Scan(&client.SyncTms)
	if err != nil {
		return domain.Client{}, fmt.Errorf("row.Scan: %w", err)
	}
	return client, nil
}

func (r *Repository) UpdateClientLastSyncTmsByID(ctx context.Context,
	id string, syncTms time.Time) error {
	query := `update keeper.client set sync_tms = @sync_tms where id = @id`
	args := pgx.NamedArgs{
		"sync_tms": syncTms,
		"id":       id,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
func (r *Repository) FindClientByID(ctx context.Context, id string) (domain.Client, error) {
	query := `select id, user_id, sync_tms from keeper.client where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.db.QueryRow(ctx, query, args)
	var client domain.Client
	err := row.Scan(&client.ID, &client.UserID, &client.SyncTms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Client{}, service.ErrItemNotFound
		}
		return domain.Client{}, fmt.Errorf("row.Scan: %w", err)
	}
	return client, nil
}
func (r *Repository) FindClientsByUserID(ctx context.Context, userID string) ([]domain.Client, error) {
	query := `select id, user_id, sync_tms from keeper.client where user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]domain.Client, 0)
	for rows.Next() {
		var c domain.Client
		errScan := rows.Scan(&c.ID, &c.UserID, &c.SyncTms)
		if errScan != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, c)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}
