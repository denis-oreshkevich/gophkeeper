package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/jackc/pgx/v5"
	"time"
)

func (r *Repository) SaveCredentials(ctx context.Context, cred model.Credentials) error {
	query := `insert into keeper.cred(id, login, password, user_id, status, modified_tms)
	values (@id, @login, @password, @user_id, @status, @modified_tms) on conflict (id) 
	do update set login = @login, password = @password, status = @status, 
	modified_tms = @modified_tms`
	args := pgx.NamedArgs{
		"id":           cred.ID,
		"login":        cred.Login,
		"password":     cred.Password,
		"user_id":      cred.UserID,
		"status":       cred.Status,
		"modified_tms": cred.ModifiedTms,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
func (r *Repository) FindCredentialsByID(ctx context.Context, id string) (model.Credentials, error) {
	query := `select id, login, password, user_id, status, modified_tms 
	from keeper.cred where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.db.QueryRow(ctx, query, args)
	var cred model.Credentials
	err := row.Scan(&cred.ID, &cred.Login, &cred.Password, &cred.UserID,
		&cred.Status, &cred.ModifiedTms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Credentials{}, repo.ErrItemNotFound
		}
		return model.Credentials{}, fmt.Errorf("row.Scan: %w", err)
	}
	return cred, nil
}

func (r *Repository) FindCredentialsByUserID(ctx context.Context,
	userID string) ([]model.Credentials, error) {
	query := `select id, login, password, user_id, status, modified_tms
	from keeper.cred where user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]model.Credentials, 0)
	for rows.Next() {
		var c model.Credentials
		errScan := rows.Scan(&c.ID, &c.Login, &c.Password, &c.UserID, &c.Status, &c.ModifiedTms)
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

func (r *Repository) FindCredentialsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]model.Credentials, error) {
	query := `select id, login, password, user_id, status, modified_tms
	from keeper.cred where user_id=@user_id and modified_tms > @tms`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]model.Credentials, 0)
	for rows.Next() {
		var c model.Credentials
		errScan := rows.Scan(&c.ID, &c.Login, &c.Password, &c.UserID, &c.Status, &c.ModifiedTms)
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

func (r *Repository) DeleteCredentialsByID(ctx context.Context, id string) error {
	query := `update keeper.card set status = @status where id = @id`
	args := pgx.NamedArgs{
		"id":     id,
		"status": model.StatusDeleted,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
