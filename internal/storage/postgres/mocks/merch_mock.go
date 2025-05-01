package mocks

import (
	"context"
	"errors"
	"merch_service/internal/models"
)

// Структура мока TODO
type MerchMock struct {

	// Добавить поля для контроля поведения мока

	forceError bool
}

// NewMerchMock создает мок
func NewMerchMock() *MerchMock {
	return &MerchMock{}
}

// ForceError метод для тестирования ошибок
func (m *MerchMock) ForceError() {
	m.forceError = true
}

// Mock реализация интерфейса MerchStorage

func (m *MerchMock) Create(ctx context.Context, merch *models.Item) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}

func (m *MerchMock) Get(ctx context.Context, name string) (*models.Item, error) {
	if m.forceError {
		return nil, errors.New("forced error")
	}
	return &models.Item{
		Name:  name,
		Price: 100,
		Stock: 10,
	}, nil
}

func (m *MerchMock) GetList(ctx context.Context) ([]*models.Item, error) {
	if m.forceError {
		return nil, errors.New("forced error")
	}
	return []*models.Item{
		{Name: "Mock Item 1", Price: 100, Stock: 5},
		{Name: "Mock Item 2", Price: 200, Stock: 3},
	}, nil
}

func (m *MerchMock) Update(ctx context.Context, name string, merch *models.Item) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}

func (m *MerchMock) Delete(ctx context.Context, name string) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}
