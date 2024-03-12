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

func (r *Repository) SaveBinary(ctx context.Context, bin *domain.Binary) error {
	query := `insert into keeper.binary(id, f_name, "data", user_id, status, modified_tms)
	values (@id, @f_name, @data, @user_id, @status, @modified_tms) on conflict do update set 
	f_name = @f_name, "data" = @data, status = @status, modified_tms = @modified_tms`
	args := pgx.NamedArgs{
		"id":           bin.ID,
		"f_name":       bin.Name,
		"data":         bin.Data,
		"user_id":      bin.UserID,
		"status":       bin.Status,
		"modified_tms": bin.ModifiedTms,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindBinaryByID(ctx context.Context, id string) (*domain.Binary, error) {
	query := `select id, f_name, "data", user_id, status, modified_tms from keeper.binary where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.db.QueryRow(ctx, query, args)
	var binary domain.Binary
	err := row.Scan(&binary.ID, &binary.Name, &binary.Data, &binary.UserID, &binary.Status, &binary.ModifiedTms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrItemNotFound
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	return &binary, nil
}

func (r *Repository) FindBinariesByUserID(ctx context.Context,
	userID string) ([]*domain.Binary, error) {
	query := `select id, f_name, "data", user_id, status, modified_tms 
	from keeper.binary where user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*domain.Binary, 0)
	for rows.Next() {
		var b domain.Binary
		errScan := rows.Scan(&b.ID, &b.Name, &b.Data, &b.UserID, &b.Status, &b.ModifiedTms)
		if errScan != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, &b)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}

func (r *Repository) FindActiveBinariesModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*domain.Binary, error) {
	query := `select id, f_name, "data", user_id, status, modified_tms 
	from keeper.binary where user_id = @user_id and modified_tms > @tms and status = @status`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
		"status":  domain.StatusActive,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*domain.Binary, 0)
	for rows.Next() {
		var b domain.Binary
		errScan := rows.Scan(&b.ID, &b.Name, &b.Data, &b.UserID, &b.Status, &b.ModifiedTms)
		if errScan != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, &b)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}

func (r *Repository) FindDeletedBinariesModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*domain.Binary, error) {
	query := `select id, user_id, status, modified_tms 
	from keeper.binary where user_id = @user_id and modified_tms > @tms and status = @status`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
		"status":  domain.StatusDeleted,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*domain.Binary, 0)
	for rows.Next() {
		var b domain.Binary
		errScan := rows.Scan(&b.ID, &b.UserID, &b.Status, &b.ModifiedTms)
		if errScan != nil {
			return nil, fmt.Errorf("cannot scan value. %w", err)
		}
		res = append(res, &b)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows.Err(). %w", err)
	}
	return res, nil
}

func (r *Repository) DeleteBinaryByID(ctx context.Context, id string) error {
	query := `update keeper.binary set status = @status where id = @id`
	args := pgx.NamedArgs{
		"id":     id,
		"status": domain.StatusDeleted,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
