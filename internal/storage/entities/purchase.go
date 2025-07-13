package entities

import (
	"context"
	"merch_service/internal/models"
)

type PurchaseStorage interface {
	// Create - добавляет покупкку мерча юзером в историю
	Create(ctx context.Context, currUser *models.User, merch *models.Item, count int) error

	// Get - получает слайс покупок пользователя
	Get(ctx context.Context, user *models.User) ([]*models.PurchaseEntry, error)
}
