package mocks

import (
	"context"
	"errors"

	"merch_service/new_version/internal/models"
)

// UserMock - моковая реализация UserStorage
type UserMock struct {

	// Добавить поля для контроля поведения мока

	forceError bool
}

// NewUserMock создает новый экземпляр UserMock
func NewUserMock() *UserMock {
	return &UserMock{}
}

// ForceError метод для тестирования ошибок
func (m *UserMock) ForceError() {
	m.forceError = true
}

//
// Mock реализация интерфейса UserStorage
//

// Create моковая реализация
func (m *UserMock) Create(ctx context.Context, user *models.User) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}

// Get моковая реализация
func (m *UserMock) Get(ctx context.Context, login string) (*models.User, error) {
	if m.forceError {
		return nil, errors.New("forced error")
	}
	return &models.User{
		Id:       1,
		Login:    login,
		Password: "mocked_password",
		Coins:    1000,
	}, nil
}

// Update моковая реализация
func (m *UserMock) Update(ctx context.Context, login string, user *models.User) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}

// Delete моковая реализация
func (m *UserMock) Delete(ctx context.Context, login string) error {
	if m.forceError {
		return errors.New("forced error")
	}
	return nil
}