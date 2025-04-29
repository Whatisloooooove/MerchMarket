package entities

import (
	"context"
	"merch_service/internal/models"
)

type CoinsStorage interface {
	// Create - добавляет изменение баланса юзера в историю
	Create(ctx context.Context, currUser *models.User, oldBalance int) error

	// Get - получает слайс изменений баланса пользователя
	Get(ctx context.Context, userName string) ([]*models.CoinsEntry, error)
}
