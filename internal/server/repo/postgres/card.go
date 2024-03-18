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

func (r *Repository) SaveCard(ctx context.Context, card model.Card) error {
	query := `insert into keeper.card(id, num, cvc, holder_name, user_id, status, modified_tms) 
	values (@id, @num, @cvc, @holder_name, @user_id, @status, @modified_tms) 
	on conflict (id) do update set num = @num, cvc = @cvc, holder_name = @holder_name, 
	status = @status, modified_tms = @modified_tms`
	args := pgx.NamedArgs{
		"id":           card.ID,
		"num":          card.Num,
		"cvc":          card.CVC,
		"holder_name":  card.HolderName,
		"user_id":      card.UserID,
		"status":       card.Status,
		"modified_tms": card.ModifiedTms,
	}
	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("db.Exec: %w", err)
	}
	return nil
}
func (r *Repository) FindCardByID(ctx context.Context, id string) (model.Card, error) {
	query := `select id, num, cvc, holder_name, user_id, status, modified_tms 
	from keeper.card where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	row := r.db.QueryRow(ctx, query, args)
	var card model.Card
	err := row.Scan(&card.ID, &card.Num,
		&card.CVC, &card.HolderName, &card.UserID, &card.Status, &card.ModifiedTms)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Card{}, repo.ErrItemNotFound
		}
		return model.Card{}, fmt.Errorf("row.Scan: %w", err)
	}
	return card, nil
}
func (r *Repository) FindCardsByUserID(ctx context.Context, userID string) ([]model.Card, error) {
	query := `select id, num, cvc, holder_name, user_id, status, modified_tms 
	from keeper.card where user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": userID,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]model.Card, 0)
	for rows.Next() {
		var c model.Card
		errScan := rows.Scan(&c.ID, &c.Num,
			&c.CVC, &c.HolderName, &c.UserID, &c.Status, &c.ModifiedTms)
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

func (r *Repository) FindCardsModifiedAfter(ctx context.Context, userID string,
	tms time.Time) ([]model.Card, error) {
	query := `select id, num, cvc, holder_name, user_id, status, modified_tms 
	from keeper.card where user_id=@user_id and modified_tms > @tms`
	args := pgx.NamedArgs{
		"user_id": userID,
		"tms":     tms,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	var res = make([]model.Card, 0)
	for rows.Next() {
		var c model.Card
		errScan := rows.Scan(&c.ID, &c.Num,
			&c.CVC, &c.HolderName, &c.UserID, &c.Status, &c.ModifiedTms)
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

func (r *Repository) DeleteCardByID(ctx context.Context, id string) error {
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
