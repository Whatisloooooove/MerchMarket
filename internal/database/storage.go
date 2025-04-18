package database

import (
	"context"
	"merch_service/internal/models"
)

type Storage interface {
	// Получение информации о пользователе по логину в виде models.User
	// Вместо нескольких запросов для получения user_id, coins и т.д.
	// можно обращатся к этой функции.
	UserByLogin(ctx context.Context, Login string) (*models.User, error)

	// UpdateUserCoins - обновляет число монет пользователя
	UpdateUserCoins(ctx context.Context, user *models.User, Amount int) error

	// CoinsHistory - возвращает историю кошелька пользователя
	CoinsHistory(ctx context.Context, user *models.User) ([]models.CoinsEntry, error)

	// Можно на уровне базы данных сделать триггер для записи
}

// -------------------------------------------------------- //
// --------- Имплементация для работы с postgres ---------- //

// Имена таблиц
const (
	schema            = "merchshop."
	usersTable        = schema + "users"
	coinHistoryTable  = schema + "coinhistory"
	merchTable        = schema + "merch"
	transactionsTable = schema + "transactions"
	purchasesTable    = schema + "purchases"
)

func (db *DB) UserByLogin(ctx context.Context, Login string) (*models.User, error) {

	var user models.User
	return &user, nil

	// requestStr := `SELECT user_id, login, password, email, coins
	// 					FROM $1 ...`
}
