package postgres

import (
	"context"
	"errors"
	"merch_service/new_version/internal/models"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// MerchPG реализует интерфейс MerchStorage в PostgreSQL
type MerchPG struct {
	db *pgxpool.Pool
}

// NewMerchStorage создает новый экземпляр хранилища мерча.
func NewMerchStorage(db *pgxpool.Pool) *MerchPG {
	return &MerchPG{db: db}
}

// validateID проверяет валидность ID товара.
func (m *MerchPG) validateID(id int) error {
	if id <= 0 {
		return models.ErrInvalidMerchID
	}
	return nil
}

// validateMerch проверяет валидность данных товара.
func (m *MerchPG) validateMerch(merch *models.Item) error {
	if merch == nil {
		return models.ErrEmptyMerch
	}
	if merch.Name == "" {
		return models.ErrEmptyMerchName
	}
	if merch.Price < 0 {
		return models.ErrNegativePrice
	}
	if merch.Stock < 0 {
		return models.ErrNegativeStock
	}
	return nil
}

// Create создает новый товар в базе данных.
// Возвращает ошибку при невалидных данных или проблемах с БД.
func (m *MerchPG) Create(ctx context.Context, merch *models.Item) error {
	if err := m.validateMerch(merch); err != nil {
		return err
	}

	query := `
		INSERT INTO merchshop.merch (name, price, stock)
		VALUES ($1, $2, $3)
		RETURNING merch_id
	`

	err := m.db.QueryRow(
		ctx,
		query,
		merch.Name,
		merch.Price,
		merch.Stock,
	).Scan(&merch.Id)

	if err != nil {
		return err
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
	}

	return nil
}

// Get возвращает товар по ID. Если товар не найден,
// возвращает nil и ошибку.
func (m *MerchPG) Get(ctx context.Context, id int) (*models.Item, error) {
	if err := m.validateID(id); err != nil {
		return nil, err
	}

	query := `
		SELECT merch_id, name, price, stock
		FROM merchshop.merch
		WHERE merch_id = $1
	`

	var item models.Item
	err := m.db.QueryRow(ctx, query, id).Scan(
		&item.Id,
		&item.Name,
		&item.Price,
		&item.Stock,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrMerchNotFound
		}
		return nil, err
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
	}

	return &item, nil
}

<<<<<<< HEAD:internal/storage/postgres/merch_pg.go
func (s *merchStorage) GetList(ctx context.Context) ([]*models.Item, error) {
	query := `
		SELECT name, price, stock 
		FROM merchshop.merch
		ORDER BY name
	`

	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query merch list: %w", err)
=======
// Update обновляет данные товара. Возвращает ошибку если:
// - товар не найден
// - данные невалидны
// - произошла ошибка БД
func (m *MerchPG) Update(ctx context.Context, id int, merch *models.Item) error {
	if err := m.validateID(id); err != nil {
		return err
	}

	if err := m.validateMerch(merch); err != nil {
		return err
	}

	query := `
		UPDATE merchshop.merch
		SET name = $1, price = $2, stock = $3
		WHERE merch_id = $4
	`

	result, err := m.db.Exec(
		ctx,
		query,
		merch.Name,
		merch.Price,
		merch.Stock,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrMerchNotFound
	}

	return nil
}

// Delete удаляет товар по ID. Возвращает ошибку если:
// - товар не найден
// - ID невалиден
// - произошла ошибка БД
func (m *MerchPG) Delete(ctx context.Context, id int) error {
	if err := m.validateID(id); err != nil {
		return err
	}

	query := `
		DELETE FROM merchshop.merch
		WHERE merch_id = $1
	`

	result, err := m.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrMerchNotFound
	}

	return nil
}

// GetList возвращает список всех товаров.
// Возвращает ошибку при проблемах с БД.
func (m *MerchPG) GetList(ctx context.Context) ([]*models.Item, error) {
	query := `
		SELECT merch_id, name, price, stock
		FROM merchshop.merch
		ORDER BY merch_id
	`

	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, err
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
<<<<<<< HEAD:internal/storage/postgres/merch_pg.go
		if err := rows.Scan(&item.Name, &item.Price, &item.Stock); err != nil {
			return nil, fmt.Errorf("failed to scan merch item: %w", err)
=======
		if err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Price,
			&item.Stock,
		); err != nil {
			return nil, err
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
<<<<<<< HEAD:internal/storage/postgres/merch_pg.go
		return nil, fmt.Errorf("rows error: %w", err)
=======
		return nil, err
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
	}

	return items, nil
}
<<<<<<< HEAD:internal/storage/postgres/merch_pg.go

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
=======
>>>>>>> develop:new_version/internal/storage/postgres/merch_pg.go
