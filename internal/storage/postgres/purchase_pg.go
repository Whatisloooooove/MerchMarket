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

// NewCoinsStorage создает новый экземпляр истории монет.
func NewPurchaseStorage(db *pgxpool.Pool) *PurchasePG {
	return &PurchasePG{db: db}
}

// Create - добавляет покупкку мерча юзером в историю
func (p *PurchasePG) Create(ctx context.Context, currUser *models.User, merch *models.Item, count int) error {
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

	return err
}

// Get - получает слайс покупок пользователя
func (p *PurchasePG) Get(ctx context.Context, user *models.User) ([]*models.PurchaseEntry, error) {
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
