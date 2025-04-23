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
