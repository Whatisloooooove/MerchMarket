package postgres

import (
	"context"

	"merch_service/new_version/internal/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

// TransactionPG реализует интерфейс TransactionStorage в PostgreSQL
type TransactionPG struct {
	db *pgxpool.Pool
}

// NewTransactionStorage создает новый экземпляр хранилища транзакций
func NewTransactionStorage(db *pgxpool.Pool) *TransactionPG {
	return &TransactionPG{db: db}
}

// validateTransaction проверяет валидность данных транзакции
func (t *TransactionPG) validateTransaction(tr *models.TransactionEntry) error {
	if tr == nil {
		return models.ErrEmptyTransaction
	}
	if tr.SenderID <= 0 {
		return models.ErrInvalidSenderID
	}
	if tr.ReceiverID <= 0 {
		return models.ErrInvalidReceiverID
	}
	if tr.Amount <= 0 {
		return models.ErrInvalidAmount
	}
	if tr.SenderID == tr.ReceiverID {
		return models.ErrSameSenderReceiver
	}
	return nil
}

// Create создает новую транзакцию между пользователями
func (t *TransactionPG) Create(ctx context.Context, tr *models.TransactionEntry) error {
	if err := t.validateTransaction(tr); err != nil {
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
		tr.SenderID,
	).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrSenderNotFound
	}

	err = tx.QueryRow(ctx, 
		"SELECT EXISTS(SELECT 1 FROM merchshop.users WHERE user_id = $1)", 
		tr.ReceiverID,
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
		tr.SenderID,
	).Scan(&senderCoins)
	if err != nil {
		return err
	}

	if senderCoins < tr.Amount {
		return models.ErrInsufficientCoins
	}

	// Выполняем перевод
	_, err = tx.Exec(ctx,
		"UPDATE merchshop.users SET coins = coins - $1 WHERE user_id = $2",
		tr.Amount, tr.SenderID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		"UPDATE merchshop.users SET coins = coins + $1 WHERE user_id = $2",
		tr.Amount, tr.ReceiverID,
	)
	if err != nil {
		return err
	}

	// Записываем транзакцию (без возвращения ID)
	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.transactions 
		(sender_id, receiver_id, amount)
		VALUES ($1, $2, $3)`,
		tr.SenderID, tr.ReceiverID, tr.Amount,
	)
	if err != nil {
		return err
	}

	// Записываем в историю изменений баланса отправителя
	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.coinhistory 
		(user_id, coins_before, coins_after)
		VALUES ($1, $2, $3)`,
		tr.SenderID, senderCoins, senderCoins-tr.Amount,
	)
	if err != nil {
		return err
	}

	// Записываем в историю изменений баланса получателя
	var receiverCoins int
	err = tx.QueryRow(ctx,
		"SELECT coins FROM merchshop.users WHERE user_id = $1",
		tr.ReceiverID,
	).Scan(&receiverCoins)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO merchshop.coinhistory 
		(user_id, coins_before, coins_after)
		VALUES ($1, $2, $3)`,
		tr.ReceiverID, receiverCoins, receiverCoins+tr.Amount,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}