package postgres

import (
	"context"

	"merch_service/internal/models"
	"merch_service/internal/storage/entities"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ entities.TransactionStorage = (*TransactionPG)(nil)

// TransactionPG реализует интерфейс TransactionStorage в PostgreSQL
type TransactionPG struct {
	db *pgxpool.Pool
}

// NewTransactionStorage создает новый экземпляр хранилища транзакций
func NewTransactionStorage(db *pgxpool.Pool) *TransactionPG {
	return &TransactionPG{db: db}
}

// validateTransaction проверяет валидность данных транзакции
func (t *TransactionPG) validateTransaction(send *models.User, recv *models.User, amount int) error {
	if send == nil || recv == nil {
		return models.ErrEmptyTransaction
	}
	if send.Id <= 0 {
		return models.ErrInvalidSenderID
	}
	if recv.Id <= 0 {
		return models.ErrInvalidReceiverID
	}
	if amount <= 0 {
		return models.ErrInvalidAmount
	}
	if send.Id == recv.Id {
		return models.ErrSameSenderReceiver
	}
	return nil
}

// validateUser проверяет валидность данных пользователя.
func (t *TransactionPG) validateUser(user *models.User) error {
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

// Create создает новую транзакцию между пользователями
func (t *TransactionPG) Create(ctx context.Context, send *models.User, recv *models.User, amount int) error {
	if err := t.validateTransaction(send, recv, amount); err != nil {
		return err
	}

	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Проверяем существование пользователей
	var exists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM merchshop.users WHERE user_id = $1)",
		send.Id,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrSenderNotFound
	}

	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM merchshop.users WHERE user_id = $1)",
		recv.Id,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrReceiverNotFound
	}

	// Проверяем баланс отправителя
	var senderCoins int
	err = tx.QueryRow(ctx,
		"SELECT coins FROM merchshop.users WHERE user_id = $1 FOR UPDATE",
		send.Id,
	).Scan(&senderCoins)
	if err != nil {
		return err
	}

	if senderCoins < amount {
		return models.ErrInsufficientCoins
	}

	// Выполняем перевод
	_, err = tx.Exec(ctx,
		"UPDATE merchshop.users SET coins = coins - $1 WHERE user_id = $2",
		amount, send.Id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		"UPDATE merchshop.users SET coins = coins + $1 WHERE user_id = $2",
		amount, recv.Id,
	)
	if err != nil {
		return err
	}

	// Записываем транзакцию (без возвращения ID)
	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.transactions 
		(sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)`,
		send.Id, recv.Id, amount,
	)
	if err != nil {
		return err
	}

	// Записываем в историю изменений баланса отправителя
	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.coinhistory 
		(user_id, coins_before, coins_after)
		VALUES ($1, $2, $3)`,
		send.Id, senderCoins, senderCoins-amount,
	)
	if err != nil {
		return err
	}

	// Записываем в историю изменений баланса получателя
	var receiverCoins int
	err = tx.QueryRow(ctx,
		"SELECT coins FROM merchshop.users WHERE user_id = $1",
		recv.Id,
	).Scan(&receiverCoins)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.coinhistory 
		(user_id, coins_before, coins_after)
		VALUES ($1, $2, $3)`,
		recv.Id, receiverCoins, receiverCoins+amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Get возвращает историю транзакций пользователя
func (t *TransactionPG) Get(ctx context.Context, user *models.User) ([]*models.TransactionEntry, error) {
	if err := t.validateUser(user); err != nil {
		return nil, err
	}

	query := `
		SELECT 
			t.transaction_id,
			t.sender_id,
			t.receiver_id,
			t.amount
		FROM merchshop.transactions t
		WHERE t.sender_id = $1 OR t.receiver_id = $1
		ORDER BY t.transaction_date DESC
	`

	rows, err := t.db.Query(ctx, query, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.TransactionEntry
	for rows.Next() {
		var entry models.TransactionEntry

		err := rows.Scan(
			&entry.Id,
			&entry.SenderID,
			&entry.ReceiverID,
			&entry.Amount,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, models.ErrNoTransactionsFound
	}

	return transactions, nil
}

// DeleteForUser удаляет все транзакции пользователя
func (t *TransactionPG) DeleteForUser(ctx context.Context, user *models.User) error {
	if err := t.validateUser(user); err != nil {
		return err
	}

	query := `DELETE FROM merchshop.transactions WHERE sender_id = $1 OR receiver_id = $1`

	result, err := t.db.Exec(ctx, query, user.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrNoTransactionsFound
	}

	return nil
}
