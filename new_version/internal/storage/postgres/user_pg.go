package postgres

import (
	"context"

	"merch_service/new_version/internal/models"

	"github.com/jackc/pgx/v4"
)

type UserPG struct {
	conn *pgx.Conn
}

func NewUserStorage(conn *pgx.Conn) *UserPG {
	return &UserPG{conn: conn}
}

func (s *UserPG) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO merchshop.users (login, password, coins) 
		VALUES ($1, $2, $3)
		ON CONFLICT (login) DO NOTHING
	`
	tag, err := s.conn.Exec(ctx, query, 
		user.Login, user.Password, user.Coins)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (s *UserPG) Get(ctx context.Context, login string) (*models.User, error) {
	query := `
		SELECT user_id, login, password, coins 
		FROM merchshop.users 
		WHERE login = $1
	`
	var user models.User
	err := s.conn.QueryRow(ctx, query, login).Scan(
		&user.Id, &user.Login, &user.Password, &user.Coins)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserPG) Update(ctx context.Context, login string, user *models.User) error {
	query := `
		UPDATE merchshop.users 
		SET login = $1, password = $2, coins = $3 
		WHERE login = $4
	`
	tag, err := s.conn.Exec(ctx, query, 
		user.Login, user.Password, user.Coins, login)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (s *UserPG) Delete(ctx context.Context, login string) error {
	query := `DELETE FROM merchshop.users WHERE login = $1`
	tag, err := s.conn.Exec(ctx, query, login)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil
	}
	return nil
}