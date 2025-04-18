package database

import (
	"context"
	"errors"
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
	requestStr := `INSERT INTO merchshop.users (login, password, email) VALUES ($1, $2, $1)`

	_, err := db.pool.Exec(context.TODO(), requestStr, user.Login, user.Password)

	if err != nil {
		log.Println("ошибка при записи в таблицу:", err)
		return err
	}
	return nil
}

func (db *DB) TransferCoins(transaction *models.TransactionEntry) error {
	// CONCURRENCY WARNING!
	// Следующие код обращения к базе данных нужно поменять, так как он
	// содает data-race!

	// Все фигня! Нужно сделать SQL функцию, чтобы делать откатыват, если произошла ошибка.
	// Через драйвер такое делать - дело неблагодарное. TODO!

	if transaction.Amount <= 0 {
		return errors.New("не смешно =(")
	}

	// TODO нужно проверять что получатель существует; Если нет,
	// то восстановить баланс (во вселенной) отправителя

	// Проверка баланса происходит на уровне базы данных
	// Если все ок, снимаем у него, и добавляем получителю
	var (
		sender_id            int
		sender_coins_after   int
		reciever_id          int
		reciever_coins_after int
		withdrawalRequest    = `UPDATE merchshop.users
			SET coins = coins - $1
			WHERE login = $2
			RETURNING user_id, coins`

		// depositRequest можно исползовать для возрвата средств отправителю,
		// если произошла ошибка
		depositRequest = `UPDATE merchshop.users
			SET coins = coins + $1
			WHERE login = $2
			RETURNING user_id, coins`

		updateTransactionsRequest = `INSERT INTO merchshop.transactions
			(sender_id, reciever_id, amount) VALUES ($1, $2, $3)`

		updateHistroyRequest = `INSERT INTO merchshop.coinhistory
			(user_id, coins_before, coins_after) VALUES ($1, $2, $3)`
		// updateRecieverHistoryRequest = '2'
	)

	err := db.pool.QueryRow(context.TODO(),
		withdrawalRequest,
		transaction.Amount,
		transaction.Sender).Scan(&sender_id, &sender_coins_after)

	if err != nil {
		log.Println("ошибка при обновлении баланса отправителя:", err)
		return err
	}

	err = db.pool.QueryRow(context.TODO(),
		depositRequest,
		transaction.Amount,
		transaction.Reciever).Scan(&reciever_id, &reciever_coins_after)

	// IMPORTANT!!!
	// Если пользователя не нашлось, UPDATE не кидает ошибку; нужна отдельная
	// функция на уровне базы данных, чтобы монеты отправителя не исчезли впустую
	if err != nil {
		log.Println("ошибка при обновлении баланса получателя:", err)
		return err
	}

	// Запись в таблицу транзакции
	_, err = db.pool.Exec(context.TODO(),
		updateTransactionsRequest,
		sender_id,
		reciever_id,
		transaction.Amount,
	)

	if err != nil {
		log.Println("ошибка при записи в транзации")
		return err
	}

	_, err = db.pool.Exec(context.TODO(),
		updateHistroyRequest,
		sender_id,
		sender_coins_after+transaction.Amount,
		sender_coins_after,
	)

	if err != nil {
		log.Println("ошибка при записи в истории отправителя")
		return err
	}

	_, err = db.pool.Exec(context.TODO(),
		updateHistroyRequest,
		reciever_id,
		reciever_coins_after-transaction.Amount,
		reciever_coins_after,
	)

	if err != nil {
		log.Println("ошибка при записи в истории получателя")
		return err
	}

	return nil
}

func (db *DB) CoinsHistory(user *models.User) ([]models.CoinsEntry, error) {
	var (
		user_id        int
		userIdRequest  = `SELECT user_id FROM merchshop.users WHERE login = $1`
		historyRequest = `SELECT coins_before, coins_after, change_date
							FROM merchshop.coinhistory
								WHERE user_id = $1;`
	)

	err := db.pool.QueryRow(context.TODO(),
		userIdRequest,
		user.Login).Scan(&user_id)

	if err != nil {
		log.Println("ошибка при посике пользователя:", err)
		return nil, err
	}

	rows, err := db.pool.Query(context.TODO(),
		historyRequest,
		user_id)

	if err != nil {
		log.Println("ошибка при чтении истории пользователя:", err)
		return nil, err
	}
	defer rows.Close()

	var history []models.CoinsEntry

	for rows.Next() {
		var entry models.CoinsEntry
		err := rows.Scan(&entry.CoinsBefore, &entry.CoinsAfter, &entry.Date)
		if err != nil {
			log.Println("ошибка чтения из базы данных:", err)
			return nil, err
		}
		history = append(history, entry)
	}

	if err := rows.Err(); err != nil {
		log.Println("ошибка чтения из базы данных:", err)
		return nil, err
	}
	return history, nil
}
