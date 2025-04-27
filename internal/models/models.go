package models

import "time"

// Эти структуры предназначены для работы с базой данных
// У серверного кода есть похожие структуры для обработки данных во входящих запросах
// (см. handlers.go:TransactionRequest,LoginRequest)

type Item struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type CoinsEntry struct {
	Id          int
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

type TransactionEntry struct {
	Id         int
	SenderID   int
	ReceiverID int
	Amount     int
}

type PurchaseEntry struct {
	Id       int
	ItemName string
	Count    int
	Date     time.Time
}

// Для хендлеров

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"pass"`
}

type PurchaseRequest struct {
	Item  string `json:"id"`
	Count int    `json:"count"`
}

type TransactionRequest struct {
	Reciever string `json:"reciever_id"`
	Amount   int    `json:"amount"`
}
