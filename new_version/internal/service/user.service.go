package service

import (
	"context"
	"fmt"
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type UserServiceInterface interface {
	// Login - проверяет есть ли такой пользователь
	// и, если был возвращет nil вместо ошибки
	Login(ctx context.Context, logReq *models.LoginRequest) error

	// Register - проверяет не было ли такого пользователя уже
	// и, если не было, добавляет его и возвращает nil
	Register(ctx context.Context, regReq *models.LoginRequest) error

	// CoinsHistory - возвращает историю изменения баланса пользователя
	CoinsHistory(ctx context.Context, login string) ([]models.CoinsEntry, error)

	// PurchaseHistory - возвращает историю покупок
	PurchaseHistory(ctx context.Context, login string) ([]models.PurchaseEntry, error)
}

var _ UserServiceInterface = (*UserService)(nil)

// UserService - реализует интерфейс UserServiceInterface
type UserService struct {
	UserStorage storage.Storage
}

// Login - проверяет есть ли такой пользователь
// и, если был возвращет nil вместо ошибки
func (u *UserService) Login(ctx context.Context, logReq *models.LoginRequest) error {
	user, err := u.UserStorage.FindUserByLogin(ctx, logReq.Login)
	if err != nil {
		return err
	}

	if user.Password != logReq.Password {
		return fmt.Errorf("неверный пароль")
	}

	return nil
}

// Register - проверяет не было ли такого пользователя уже
// и, если не было, добавляет его и возвращает nil
func (u *UserService) Register(ctx context.Context, regReq *models.LoginRequest) error {
	_, err := u.UserStorage.FindUserByLogin(ctx, regReq.Login)
	if err == nil {
		return fmt.Errorf("такой пользователь уже существует")
	}

	err = u.UserStorage.CreateUser(ctx, &models.User{Login: regReq.Login, Password: regReq.Password})
	if err != nil {
		return err
	}
	return nil
}

// CoinsHistory - проверяет существует ли переданный пользователь
// и возвращает слайс с историей изменения баланса
func (u *UserService) CoinsHistory(ctx context.Context, login string) ([]models.CoinsEntry, error) {
	user, err := u.UserStorage.FindUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	coinsHistory, err := u.UserStorage.CoinsHistory(ctx, user)
	if err != nil {
		return nil, err
	}
	return coinsHistory, nil
}

// PurchaseHistory - проверяет существует ли переданный пользователь
// и возвращает слайс с историей покупок мерча
func (u *UserService) PurchaseHistory(ctx context.Context, login string) ([]models.PurchaseEntry, error) {
	user, err := u.UserStorage.FindUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	purchaseHistory, err := u.UserStorage.PurchaseHistory(ctx, user)
	if err != nil {
		return nil, err
	}
	return purchaseHistory, nil
}
