package storage

import (
	"context"
	"merch_service/new_version/internal/models"
)

type Storage interface {
	// FindUserByLogin - ищет по имени пользователя и возвращает структуру User и ошибку
	FindUserByLogin(ctx context.Context, login string) (*models.User, error)

	// MakeTransaction - добавляет перевод в историю
	MakeTransaction(ctx context.Context, send *models.User, recv *models.User, amount int) error

	// CoinsHistory - возвращает историю кошелька пользователя
	CoinsHistory(ctx context.Context, user *models.User) ([]models.CoinsEntry, error)

	// PurchaseHistory - возвращает историю покупок
	PurchaseHistory(ctx context.Context, user *models.User) ([]models.PurchaseEntry, error)

	// MerchList - возвращает слайс мерчей и ошибку
	MerchList(ctx context.Context) ([]models.Item, error)

	// MerchByName -  возвращает по имени структуру Item и ошибку
	MerchByName(ctx context.Context, merchName string) (models.Item, error)

	// CreateUser - добавляет юзера в базу данных
	CreateUser(ctx context.Context, user *models.User) error

	// Buy - покупает мерч и возвращает баланс
	Buy(ctx context.Context, user *models.User, merchName string, count int) (int, error)
}
