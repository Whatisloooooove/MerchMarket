package entities

import (
	"context"

	"merch_service/internal/models"
)

// CoinsStorage определяет контракт для работы с товарами
type CoinsStorage interface {
	// Базовые CRUD операции

	// Create - добавляет изменение баланса юзера в историю.
	// Принимает экземпляр User с новым балансом и его старый баланс.
	// В случае неудачи возвращает ошибку
	Create(ctx context.Context, currUser *models.User, oldBalance int) error

	// Get - возвращает слайс изменений баланса конкретного пользователя
	// В случае неудачи возвращает ошибку
	Get(ctx context.Context, user *models.User) ([]*models.CoinsEntry, error)

	// Delete - удаляет все записи об изменении балланса экземляра User
	Delete(ctx context.Context, user *models.User) error
}
