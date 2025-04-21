package entites

import (
	"context"
	"merch_service/new_version/internal/models"
)

// MerchStorage определяет контракт для работы с товарами
type MerchStorage interface {
	// Create создает новый экземпляр мерча на основе экземпляра Item, возвращает ошибку
	Create(ctx context.Context, merch *models.Item) error

	// GetByName возвращает по name экземпляр мерча и ошибку
	Get(ctx context.Context, name string) (*models.Item, error)

	// GetList возвращает слайс всех мерчей и ошибку
	GetList(ctx context.Context) ([]*models.Item, error)

	// Update обновляет информацию по name на основе экземпляра Item, возвращает ошибку
	Update(ctx context.Context, name string, merch *models.Item) error

	// Delete удаляет мерч по name, возвращает ошибку
	Delete(ctx context.Context, name string) error
}
