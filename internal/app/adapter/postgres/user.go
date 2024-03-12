package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain"
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) CreateUser(ctx context.Context, usr domain.User) (domain.User, error) {
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
					pgErr.ConstraintName, service.ErrUserAlreadyExist)
			}
		}
		return usr, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}
func (r *Repository) FindUserByLogin(ctx context.Context, login string) (domain.User, error) {
	query := `select id, login, password from keeper.usr where login=@login`
	args := pgx.NamedArgs{
		"login": login,
	}
	row := r.db.QueryRow(ctx, query, args)
	var usr domain.User
	err := row.Scan(&usr.ID, &usr.Login, &usr.HashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, service.ErrItemNotFound
		}
		return domain.User{}, fmt.Errorf("row.Scan: %w", err)
	}
	return usr, nil
}
