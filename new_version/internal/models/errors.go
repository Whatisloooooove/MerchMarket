package models

import "errors"

// Для MerchService
var (
	ErrNoMerchInStock = errors.New("товара нет в наличии")
	ErrNoSuchMerch    = errors.New("товара по запросу не существует")
	ErrNotEnoughMerch = errors.New("на складе нет столько товара")
	ErrNotEnoughCoins = errors.New("недостаточно монет для покупки")
)

// Для UserService
var (
	ErrWrongPassword = errors.New("неверный пароль")
	ErrUserExists    = errors.New("пользователь с таким логином уже существует")
)

//
// StorageErrorsBlock
//

// Для UserStorage
var (
	ErrUserNotFound  = errors.New("такого пользователя нет в бд")
	ErrInvalidUserID = errors.New("id пользователя не может быть отрицательным")
	ErrEmptyUser     = errors.New("пользователь не может быть nill")
	ErrEmptyLogin    = errors.New("логин пользователя не может быть пустым")
	ErrEmptyPassword = errors.New("пароль пользователя не может быть пустым")
	ErrNegativeCoins = errors.New("количество монет пользователя не может быть отрицательным")
)
