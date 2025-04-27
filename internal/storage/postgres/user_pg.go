package postgres

import (
	"context"
	"errors"

	"merch_service/internal/models"
	"merch_service/internal/storage/entities"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ entities.UserStorage = (*UserPG)(nil)

// UserPG реализует интерфейс UserStorage в PostgreSQL
type UserPG struct {
	db *pgxpool.Pool
}

// NewUserStorage создает новый экземпляр хранилища пользователей.
func NewUserStorage(db *pgxpool.Pool) *UserPG {
	return &UserPG{db: db}
}

// validateID проверяет валидность ID пользователя.
func (u *UserPG) validateID(id int) error {
	if id <= 0 {
		return models.ErrInvalidUserID
	}
	return nil
}

// validateLogin проверяет валидность логина пользователя.
func (u *UserPG) validateLogin(login string) error {
	if login == "" {
		return models.ErrEmptyUserLogin
	}
	return nil
}

// validateUser проверяет валидность данных пользователя.
func (u *UserPG) validateUser(user *models.User) error {
	if user == nil {
		return models.ErrEmptyUser
	}
	if user.Login == "" {
		return models.ErrEmptyUserLogin
	}
	if user.Password == "" {
		return models.ErrEmptyUserPassword
	}
	if user.Coins < 0 {
		return models.ErrNegativeUserCoins
	}
	return nil
}

// Create создает нового пользователя в базе данных.
// Возвращает ошибку при невалидных данных или проблемах с БД.
func (u *UserPG) Create(ctx context.Context, user *models.User) error {
	if err := u.validateUser(user); err != nil {
		return err
	}

	query := `
		INSERT INTO merchshop.users (login, password, coins)
		VALUES ($1, $2, $3)
		RETURNING user_id
	`

	err := u.db.QueryRow(
		ctx,
		query,
		user.Login,
		user.Password,
		user.Coins,
	).Scan(&user.Id)

	if err != nil {
		return err
	}

	return nil
}

// Get возвращает пользователя по ID. Если пользователь не найден,
// возвращает nil и ошибку.
func (u *UserPG) Get(ctx context.Context, id int) (*models.User, error) {
	if err := u.validateID(id); err != nil {
		return nil, err
	}

	query := `
		SELECT user_id, login, password, coins
		FROM merchshop.users
		WHERE user_id = $1
	`

	var user models.User
	err := u.db.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Coins,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update обновляет данные пользователя. Возвращает ошибку если:
// - пользователь не найден
// - данные невалидны
// - произошла ошибка БД
func (u *UserPG) Update(ctx context.Context, user *models.User) error {
	if err := u.validateID(user.Id); err != nil {
		return err
	}

	if err := u.validateUser(user); err != nil {
		return err
	}

	query := `
		UPDATE merchshop.users
		SET login = $1, password = $2, coins = $3
		WHERE user_id = $4
	`

	result, err := u.db.Exec(
		ctx,
		query,
		user.Login,
		user.Password,
		user.Coins,
		user.Id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

// Delete удаляет пользователя по ID. Возвращает ошибку если:
// - пользователь не найден
// - ID невалиден
// - произошла ошибка БД
func (u *UserPG) Delete(ctx context.Context, id int) error {
	if err := u.validateID(id); err != nil {
		return err
	}

	query := `
		DELETE FROM merchshop.users
		WHERE user_id = $1
	`

	result, err := u.db.Exec(ctx, query, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

// GetByLogin возвращает пользователя по логину. Если пользователь не найден,
// возвращает nil и ошибку.
func (u *UserPG) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	if err := u.validateLogin(login); err != nil {
		return nil, err
	}

	query := `
		SELECT user_id, login, password, coins
		FROM merchshop.users
		WHERE login = $1
	`

	var user models.User
	err := u.db.QueryRow(ctx, query, login).Scan(
		&user.Id,
		&user.Login,
		&user.Password,
		&user.Coins,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserPG) GetCoinsHistory(ctx context.Context, userId int) ([]models.CoinsEntry, error) {
	// TODO
	return nil, nil
}

func (u *UserPG) GetPurchaseHistory(ctx context.Context, userId int) ([]models.PurchaseEntry, error) {
	// TODO
	return nil, nil
}
