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

func (r *Repository) SaveText(ctx context.Context, txt *model.Text) error {
	query := `insert into keeper.txt(id, val, user_id, status, modified_tms)
	values (@id, @txt, @user_id, @status, @modified_tms) on conflict (id) do update set 
	val = @txt, status = @status, modified_tms = @modified_tms`
	args := pgx.NamedArgs{
		"id":           txt.ID,
		"txt":          txt.Txt,
		"user_id":      txt.UserID,
		"status":       txt.Status,
		"modified_tms": txt.ModifiedTms,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil

}
func (r *Repository) FindTextByID(ctx context.Context, id string) (*model.Text, error) {
	query := `select id, val, user_id, status, modified_tms from keeper.txt where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.db.QueryRow(ctx, query, args)
	var txt model.Text
	err := row.Scan(&txt.ID, &txt.Txt, &txt.UserID, &txt.Status, &txt.ModifiedTms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repo.ErrItemNotFound
		}
		return nil, fmt.Errorf("row.Scan: %w", err)
	}
	return &txt, nil
}
func (r *Repository) FindTextsByUserID(ctx context.Context, userID string) ([]*model.Text, error) {
	query := `select id, val, user_id, status, modified_tms from keeper.txt where user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*model.Text, 0)
	for rows.Next() {
		var b model.Text
		errScan := rows.Scan(&b.ID, &b.Txt, &b.UserID, &b.Status, &b.ModifiedTms)
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

func (r *Repository) FindActiveTextsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Text, error) {
	query := `select id, val, user_id, status, modified_tms from keeper.txt 
    where user_id=@user_id and modified_tms > @tms and status = @status`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
		"status":  model.StatusActive,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*model.Text, 0)
	for rows.Next() {
		var b model.Text
		errScan := rows.Scan(&b.ID, &b.Txt, &b.UserID, &b.Status, &b.ModifiedTms)
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

func (r *Repository) FindDeletedTextsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]*model.Text, error) {
	query := `select id, user_id, status, modified_tms from keeper.txt 
    where user_id=@user_id and modified_tms > @tms and status = @status`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
		"status":  model.StatusDeleted,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]*model.Text, 0)
	for rows.Next() {
		var b model.Text
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

func (r *Repository) DeleteTextByID(ctx context.Context, id string) error {
	query := `update keeper.bin set status = @status where id = @id`
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
