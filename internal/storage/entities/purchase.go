package entities

import (
	"context"
	"merch_service/internal/models"
)

type PurchaseStorage interface {
	// Базовые CRUD операции

	// Create - добавляет покупку count экземляров мерча юзером. 
	// Из экземляров User и Merch записывает id
	// В случае неудачи возвращает ошибку
	Create(ctx context.Context, currUser *models.User, merch *models.Item, count int) error

	// Get - возвращает слайс покупок пользователя. 
	// В случае неудачи возвращает ошибку
	Get(ctx context.Context, user *models.User) ([]*models.PurchaseEntry, error)

	// Delete удаляет записb о покупке пользоватпеля.
    // В случае неудачи возвращает ошибку
    Delete(ctx context.Context, user *models.User) error
}
