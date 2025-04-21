package postgres

import (
	"context"
	"errors"
	"fmt"
	"merch_service/new_version/internal/models"

	"github.com/jackc/pgx/v4"
)

type merchStorage struct {
	conn *pgx.Conn
}

func NewMerchStorage(conn *pgx.Conn) *merchStorage {
	return &merchStorage{conn: conn}
}

func (s *merchStorage) Create(ctx context.Context, merch *models.Item) error {
	query := `
		INSERT INTO merchshop.merch (name, price, stock) 
		VALUES ($1, $2, $3)
		ON CONFLICT (name) DO NOTHING
	`

	tag, err := s.conn.Exec(ctx, query, merch.Name, merch.Price, merch.Stock)
	if err != nil {
		return fmt.Errorf("failed to create merch: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("merch with this name already exists")
	}

	return nil
}

func (s *merchStorage) Get(ctx context.Context, name string) (*models.Item, error) {
	query := `
		SELECT name, price, stock 
		FROM merchshop.merch 
		WHERE name = $1
	`

	var item models.Item
	err := s.conn.QueryRow(ctx, query, name).
		Scan(&item.Name, &item.Price, &item.Stock)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get merch: %w", err)
	}

	return &item, nil
}

func (s *merchStorage) GetList(ctx context.Context) ([]*models.Item, error) {
	query := `
		SELECT name, price, stock 
		FROM merchshop.merch
		ORDER BY name
	`

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query merch list: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.Name, &item.Price, &item.Stock); err != nil {
			return nil, fmt.Errorf("failed to scan merch item: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return items, nil
}

func (s *merchStorage) Update(ctx context.Context, name string, merch *models.Item) error {
	query := `
		UPDATE merchshop.merch 
		SET name = $1, price = $2, stock = $3 
		WHERE name = $4
	`

	tag, err := s.conn.Exec(ctx, query, merch.Name, merch.Price, merch.Stock, name)
	if err != nil {
		return fmt.Errorf("failed to update merch: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return nil
	}

	return nil
}

func (s *merchStorage) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM merchshop.merch WHERE name = $1`

	tag, err := s.conn.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete merch: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return nil
	}

	return nil
}
