package entites

import (
	"context"

	"merch_service/new_version/internal/models"
)

// UserStorage определяет контракт для работы с пользователями.
type UserStorage interface {
	// Create создает новый экземпляр пользователя в бд на основе экземпляра User
	// Возвращает ошибку, если создание не удалось
	Create(ctx context.Context, user *models.User) error

	// GetByLogin возвращает пользователя по его логину и ошибку.
	Get(ctx context.Context, login string) (*models.User, error)

	// Update обновляет информацию о пользователе по логину на основе экземпляра User, возвращает ошибку
	Update(ctx context.Context, login string, user *models.User) error

	// Delete удаляет пользователя по его логину. Возвращает ошибку
	Delete(ctx context.Context, login string) error
}
