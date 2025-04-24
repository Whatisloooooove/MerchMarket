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
	ErrUserNotFound      = errors.New("такого пользователя нет в бд")
	ErrEmptyUser         = errors.New("пользователь не может быть nill")
	ErrInvalidUserID     = errors.New("id пользователя не может быть отрицательным")
	ErrEmptyUserLogin    = errors.New("логин пользователя не может быть пустым")
	ErrEmptyUserPassword = errors.New("пароль пользователя не может быть пустым")
	ErrNegativeUserCoins = errors.New("количество монет пользователя не может быть отрицательным")
)

// Для MerchStorage
var (
	ErrMerchNotFound     = errors.New("такого мерча нет в бд")
	ErrEmptyMerch        = errors.New("мерч не может быть nill")
	ErrInvalidMerchID    = errors.New("id мерча не может быть отрицательным")
	ErrEmptyMerchName    = errors.New("имя мерча не может быть пустым")
	ErrNegativePrice     = errors.New("цена мерча не может быть отрицательной")
	ErrNegativeStock     = errors.New("количество мерча не может быть отрицательным")
)
