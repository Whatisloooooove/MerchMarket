package service

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entities"
)

type MerchServiceInterface interface {
	// Buy проверяет есть ли мерч в наличии и хватает
	// ли денег пользователю, после чего сохраняет данные
	// и возвращает текущий баланс
	Buy(ctx context.Context, userName, merchName string, count int) (int, error)

	// MerchList - возвращает весь доступный для покупки мерч
	MerchList(ctx context.Context) ([]*models.Item, error)
}

var _ MerchServiceInterface = (*MerchService)(nil)

// MerchService - реализует интерфейс MerchServiceInterface
type MerchService struct {
	MerchStorage entities.MerchStorage
	UserStorage  entities.UserStorage
}

// NewMerchService - создает объект MerchService
func NewMerchService(m entities.MerchStorage, u entities.UserStorage) *MerchService {
	return &MerchService{
		MerchStorage: m,
		UserStorage:  u,
	}
}

// Buy - проверяет наличие мерча и возможность пользователя купить мерч и
//
//	далее совершает покупку мерча
func (m *MerchService) Buy(ctx context.Context, userName, merchName string, count int) (int, error) {
	merch, err := m.MerchStorage.GetByName(ctx, merchName)
	if err != nil {
		return -1, err
	}

	if merch.Stock < count {
		return -1, models.ErrNotEnoughMerch
	}

	user, err := m.UserStorage.GetByLogin(ctx, userName)
	if err != nil {
		return -1, err
	}

	if user.Coins < merch.Price*count {
		return -1, models.ErrNotEnoughCoins
	}

	user.Coins -= merch.Price * count
	merch.Stock -= count

	err = m.MerchStorage.Update(ctx, merch)
	if err != nil {
		user.Coins += merch.Price * count
		merch.Stock += count
		return -1, err
	}

	err = m.UserStorage.Update(ctx, user)
	if err != nil {
		user.Coins += merch.Price * count
		merch.Stock += count
		return -1, err
	}

	return user.Coins, nil
}

// MerchList - пробрасывает контекст ниже и ждёт слайс мерчей, чтобы вернуть его
func (m *MerchService) MerchList(ctx context.Context) ([]*models.Item, error) {
	return m.MerchStorage.GetList(ctx)
}
