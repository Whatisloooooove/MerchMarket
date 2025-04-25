package entites

import (
	"context"
	"merch_service/internal/models"
)

// TransactionStorage определяет контракт для работы с товарами
type TransactionStorage interface {
	// Create - добавляет перевод в историю
	Create(ctx context.Context, send *models.User, recv *models.User, amount int) error
}
