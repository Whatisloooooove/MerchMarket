package entities

import (
	"context"

	"merch_service/internal/models"
)

// TransactionStorage определяет контракт для работы с транзакциями
type TransactionStorage interface {
	// Базовые CRUD операции

	// Create создает новую транзакцию между пользователями.
	// Возвращает ошибку при неудаче.
	Create(ctx context.Context, send *models.User, recv *models.User, amount int) error

	// Get возвращает историю транзакций пользователя
	// Возвращает ошибку при неудаче.
	Get(ctx context.Context, user *models.User) ([]*models.TransactionEntry, error)

	// DeleteForUser удаляет все транзакции пользователя
	// Возвращает ошибку при неудаче.
	DeleteForUser(ctx context.Context, user *models.User) error
}
