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
}
