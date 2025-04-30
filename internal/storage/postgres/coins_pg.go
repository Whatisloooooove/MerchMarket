package postgres

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ entities.CoinsStorage = (*CoinsPG)(nil)

// CoinsPG реализует интерфейс CoinsStroage в PostgreSQL
type CoinsPG struct {
	db *pgxpool.Pool
}

// NewCoinsStorage создает новый экземпляр истории монет.
func NewCoinsStorage(db *pgxpool.Pool) *CoinsPG {
	return &CoinsPG{db: db}
}

// validateUser проверяет валидность данных пользователя.
func (c *CoinsPG) validateUser(user *models.User) error {
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

func (c *CoinsPG) validateCoins(coins int) error {
	if coins < 0 {
		return models.ErrNegativeUserCoins
	}
	return nil
}

// Create - добавляет изменение баланса юзера в историю
func (c *CoinsPG) Create(ctx context.Context, currUser *models.User, oldBalance int) error {
	if err := c.validateUser(currUser); err != nil {
		return err
	}

	if err := c.validateCoins(oldBalance); err != nil {
		return err
	}

	query := `
		INSERT INTO merchshop.coinhistory (user_id, coins_before, coins_after)
		VALUES ($1, $2, $3);
	`

	_, err := c.db.Exec(
		ctx,
		query,
		currUser.Id,
		oldBalance,
		currUser.Coins,
	)

	if err != nil {
		return err
	}

	return nil
}

// Get - получает слайс изменений баланса пользователя
func (c *CoinsPG) Get(ctx context.Context, user *models.User) ([]*models.CoinsEntry, error) {
	if err := c.validateUser(user); err != nil {
		return nil, err
	}

	query := `
		SELECT change_id, change_date, coins_before, coins_after
		FROM merchshop.coinhistory
		WHERE user_id = $1;
		ORDER BY change_date DESC
	`

	rows, err := c.db.Query(
		ctx,
		query,
		user.Id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var coinsHist []*models.CoinsEntry
	for rows.Next() {
		var entry models.CoinsEntry
		if err := rows.Scan(
			&entry.Id,
			&entry.Date,
			&entry.CoinsBefore,
			&entry.CoinsAfter,
		); err != nil {
			return nil, err
		}
		coinsHist = append(coinsHist, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return coinsHist, nil
}

// Delete удаляет всю историю изменений баланса пользователя
func (c *CoinsPG) Delete(ctx context.Context, user *models.User) error {
	if err := c.validateUser(user); err != nil {
		return err
	}

	query := `DELETE FROM merchshop.coinhistory WHERE user_id = $1`

	result, err := c.db.Exec(ctx, query, user.Id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNoHistoryUserFound
	}

	return nil
}
