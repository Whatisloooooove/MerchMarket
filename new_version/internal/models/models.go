package models

import "time"

// Эти структуры предназначены для работы с базой данных
// У серверного кода есть похожие структуры для обработки данных во входящих запросах
// (см. handlers.go:TransactionRequest,LoginRequest)

type Item struct {
	Id int    `json:"id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}

type CoinsEntry struct {
	Date        time.Time `json:"change_date"`
	CoinsBefore int       `json:"coins_before"`
	CoinsAfter  int       `json:"coins_after"`
}

type User struct {
	Id       int
	Coins    int
	Login    string
	Password string
}

type LoginRequest struct {
	Login    string
	Password string
}

type TransactionEntry struct {
	Sender   string
	Reciever string
	Amount   int
}

type PurchaseEntry struct {
	ItemName string
	Count    int
	Date     time.Time
}
