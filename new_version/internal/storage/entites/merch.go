package entites

import (
	"context"
	"merch_service/new_version/internal/models"
)

// MerchStorage определяет контракт для работы с товарами
type MerchStorage interface {
	// Create создает новый экземпляр мерча, возвращает ошибку
	Create(ctx context.Context, merch *models.Item) error

	// GetByName возвращает по имени экземпляр мерча и ошибку
	GetByName(ctx context.Context, name string) (*models.Item, error)

	// GetList возвращает слайс всех мерчей и ошибку
	GetList(ctx context.Context) ([]*models.Item, error)

	// Update обновляет информацию о мерче, возвращает ошибку
	Update(ctx context.Context, merch *models.Item) error

	// Delete удаляет мерч по имени, возвращает ошибку
	Delete(ctx context.Context, name string) error
}
