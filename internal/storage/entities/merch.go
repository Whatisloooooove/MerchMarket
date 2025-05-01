package entities

import (
	"context"

	"merch_service/internal/models"
)

// MerchStorage определяет контракт для работы с товарами
type MerchStorage interface {
	// Базовые CRUD операции

	// Create создает новый экземпляр мерча на основе экземпляра Item,
	// обновляет ID экземпляра.
	// Возвращает ошибку при неудаче.
	Create(ctx context.Context, merch *models.Item) error

	// Get возвращает мерч по ID. Если мерч не найден,
	// возвращает nil и ошибку.
	Get(ctx context.Context, id int) (*models.Item, error)

	// Get возвращает мерч по name. Если мерч не найден,
	// возвращает nil и ошибку.
	GetByName(ctx context.Context, merchName string) (*models.Item, error)

	// Update обновляет данные в БД на основе полей экземпляра Item.
	// Возвращает ошибку при неудаче.
	Update(ctx context.Context, merch *models.Item) error

	// Delete удаляет мерч в БД по ID.
	// Возвращает ошибку при неудаче.
	Delete(ctx context.Context, id int) error

	// Дополнительные методы

	// GetList возвращает слайс всех мерчей.
	// Возвращает ошибку при неудаче.
	GetList(ctx context.Context) ([]*models.Item, error)
}
