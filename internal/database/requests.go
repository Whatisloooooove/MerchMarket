package database

import (
	"context"
	"log"
	"merch_service/internal/models"
)

// GetMerchList - возвращает список товаров для дальнейей обработки в хендлере server.MerchList
func (db *DB) GetMerchList() ([]models.Item, error) {
	requestStr := "SELECT name, price, stock FROM merchshop.merch;"
	rows, err := db.pool.Query(context.TODO(), requestStr)

	if err != nil {
		// TODO
		return nil, err
	}
	// Нужно ли?
	defer rows.Close()

	var items []models.Item

	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.Name, &item.Price, &item.Stock)
		if err != nil {
			log.Println("ошибка чтения из базы данных:", err)
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("ошибка чтения из базы данных:", err)
		return nil, err
	}
	return items, nil
}

func (db *DB) CheckUser(user *models.User) (bool, error) {
	requestStr := `SELECT EXISTS
						(SELECT FROM merchshop.users
							WHERE login = $1 AND password = $2);`
	row := db.pool.QueryRow(context.TODO(), requestStr, user.Login, user.Password)

	var registered bool
	err := row.Scan(&registered)

	if err != nil {
		log.Println("ошибка проверки регистрации:", err)
		return false, err
	}

	return registered, nil
}

func (db *DB) RegisterUser(user *models.User) error {
	// запись email TODO
	requestStr := `INSERT INTO merchshop.users (login, password, email) VALUES ($1, $2, 'tmp@mail.ru');`

	_, err := db.pool.Exec(context.TODO(), requestStr, user.Login, user.Password)

	if err != nil {
		log.Println("ошибка при записи в таблицу:", err)
		return err
	}
	return nil
}
