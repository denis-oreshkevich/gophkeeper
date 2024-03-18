package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/server/repo"
	"github.com/denis-oreshkevich/gophkeeper/internal/shared/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, usr model.User) (model.User, error) {
	query := `insert into keeper.usr(login, password) values (@login, @password) returning usr.id`
	args := pgx.NamedArgs{
		"login":    usr.Login,
		"password": usr.HashedPassword,
	}
	row := r.db.QueryRow(ctx, query, args)
	err := row.Scan(&usr.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.UniqueViolation == pgErr.Code {
				return usr, fmt.Errorf("row.Scan, contraint %s: %w",
					pgErr.ConstraintName, repo.ErrUserAlreadyExist)
			}
		}
		return usr, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}
func (r *Repository) FindUserByLogin(ctx context.Context, login string) (model.User, error) {
	query := `select id, login, password from keeper.usr where login=@login`
	args := pgx.NamedArgs{
		"login": login,
	}
	row := r.db.QueryRow(ctx, query, args)
	var usr model.User
	err := row.Scan(&usr.ID, &usr.Login, &usr.HashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repo.ErrItemNotFound
		}
		return model.User{}, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}
