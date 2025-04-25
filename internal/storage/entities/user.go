package entities

import (
	"context"

	"merch_service/internal/models"
)

// UserStorage определяет интерфейс для операций с пользователями в хранилище.
type UserStorage interface {
	// Базовые CRUD операции

	// Create создает нового пользователя в БД
	// на основе экземпляра User и обновляет ID в экземпляре.
	// Возвращает ошибку при неудаче.
	Create(ctx context.Context, user *models.User) error

	// Get возвращает пользователя по ID. Если пользователь не найден,
	// возвращает nil и ошибку.
	Get(ctx context.Context, id int) (*models.User, error)

	// Update обновляет данные в БД с id на основе полей экземпляра User.
	// Возвращает ошибку при неудаче.
	Update(ctx context.Context, user *models.User) error

	// Delete удаляет пользователя в БД по ID.
	// Возвращает ошибку при неудаче.
	Delete(ctx context.Context, id int) error

	// Дополнительные методы

	// Get возвращает пользователя по ID. Если пользователь не найден,
	// возвращает nil и ошибку.
	GetByLogin(ctx context.Context, login string) (*models.User, error)

	// GetCoinsHistory - возвращает историю кошелька пользователя
	GetCoinsHistory(ctx context.Context, userId int) ([]models.CoinsEntry, error)

	// GetPurchaseHistory - возвращает историю покупок
	GetPurchaseHistory(ctx context.Context, userId int) ([]models.PurchaseEntry, error)
}
