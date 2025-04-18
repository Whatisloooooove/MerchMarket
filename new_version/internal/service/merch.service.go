package service

import (
	"context"
	"fmt"
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type MerchServiceInterface interface {
	// Buy проверяет есть ли мерч в наличии и хватает
	// ли денег пользователю, после чего сохраняет данные
	// и возвращает текущий баланс
	Buy(ctx context.Context, login, merchName string, count int) (int, error)

	// MerchList - возвращает весь доступный для покупки мерч
	MerchList(ctx context.Context) ([]models.Item, error)
}

var _ MerchServiceInterface = (*MerchService)(nil)

// MerchService - реализует интерфейс MerchServiceInterface
type MerchService struct {
	MerchStorage storage.Storage
}

// Buy - проверяет наличие мерча и возможность пользователя купить мерч и
//
//	далее совершает покупку мерча
func (m *MerchService) Buy(ctx context.Context, login, merchName string, count int) (int, error) {
	merch, err := m.MerchStorage.MerchByName(ctx, merchName)
	if err != nil {
		return -1, err
	}

	if merch.Stock < count {
		return -1, fmt.Errorf("на складе нет такого количества мерча")
	}

	user, err := m.MerchStorage.FindUserByLogin(ctx, login)
	if err != nil {
		return -1, err
	}

	if user.Coins < merch.Price*count {
		return -1, fmt.Errorf("у вас недостаточно средств")
	}

	balance, err := m.MerchStorage.Buy(ctx, user, merchName, count)
	if err != nil {
		return -1, err
	}

	return balance, nil
}

// MerchList - пробрасывает контекст ниже и ждёт слайс мерчей, чтобы вернуть его
func (m *MerchService) MerchList(ctx context.Context) ([]models.Item, error) {
	return m.MerchStorage.MerchList(ctx)
}
