package postgres

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ entities.PurchaseStorage = (*PurchasePG)(nil)

// CoinsPG реализует интерфейс CoinsStroage в PostgreSQL
type PurchasePG struct {
	db *pgxpool.Pool
}

// validateUser проверяет валидность данных пользователя.
func (p *PurchasePG) validateUser(user *models.User) error {
	if user == nil {
		return models.ErrEmptyUser
	}
	if user.Login == "" {
		return models.ErrEmptyUserLogin
	}
	if user.Password == "" {
		return models.ErrEmptyUserPassword
	}
	if user.Coins < 0 {
		return models.ErrNegativeUserCoins
	}
	return nil
}

// validateMerch проверяет валидность данных товара.
func (p *PurchasePG) validateMerch(merch *models.Item) error {
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

// validatePurchase проверяет валидность данных покупки
func (p *PurchasePG) validatePurchase(count int) error {
	if count <= 0 {
		return models.ErrNegativeStock
	}
	return nil
}

// NewCoinsStorage создает новый экземпляр истории монет.
func NewPurchaseStorage(db *pgxpool.Pool) *PurchasePG {
	return &PurchasePG{db: db}
}

// Create - добавляет покупкку мерча юзером в историю
func (p *PurchasePG) Create(ctx context.Context, currUser *models.User, merch *models.Item, count int) error {
	if err := p.validatePurchase(count); err != nil {
		return err
	}
	
	if err := p.validateUser(currUser); err != nil {
		return err
	}
	
	if err := p.validateMerch(merch); err != nil {
		return err
	}

	query := `
		INSERT INTO merchshop.purchases (user_id, merch_id, count)
		VALUES ($1, $2, $3);
	`

	_, err := p.db.Exec(
		ctx,
		query,
		currUser.Id,
		merch.Id,
		count,
	)

	if err != nil {
		return err
	}

	return nil
}

// Get - получает слайс покупок пользователя
func (p *PurchasePG) Get(ctx context.Context, user *models.User) ([]*models.PurchaseEntry, error) {
	if err := p.validateUser(user); err != nil {
		return nil, err
	}

	query := `
		SELECT
			p.purchase_id,
			m.name,
			p.count,
			p.purchase_date
		FROM merchshop.purchases AS p
		JOIN merchshop.merch AS m  ON p.merch_id = m.merch_id
		WHERE p.user_id = $1; 
	`

	rows, err := p.db.Query(
		ctx,
		query,
		user.Id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var purchHist []*models.PurchaseEntry
	for rows.Next() {
		var entry models.PurchaseEntry
		if err := rows.Scan(
			&entry.Id,
			&entry.ItemName,
			&entry.Count,
			&entry.Date,
		); err != nil {
			return nil, err
		}
		purchHist = append(purchHist, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return purchHist, nil
}

// Delete удаляет все записи о покупках пользователя
func (p *PurchasePG) Delete(ctx context.Context, user *models.User) error {
	if err := p.validateUser(user); err != nil {
		return err
	}

	query := `
		DELETE FROM merchshop.purchases
		WHERE user_id = $1
	`

	result, err := p.db.Exec(ctx, query, user.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNoPurhasesUserFound
	}

	return nil
}
