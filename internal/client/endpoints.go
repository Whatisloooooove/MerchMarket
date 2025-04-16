package client

import "merch_service/internal/models"

// [source](https://medium.com/@bojanmajed/standard-json-api-response-format-c6c1aabcaa6d)

// ResponseBody структура для хранения ответа сервера
// см. контракт описанный выше
type ResponseBody struct {
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}

// UserTokens структура ответа на запрос к /login
type UserTokens struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

// MerchList структура ответа на запрос к /merch
type MerchList struct {
	Items []models.Item `json:"items"`
}

// HistoryLog структура ответа на запрос к /history
type HistoryLog struct {
	History []models.CoinsEntry `json:"history"`
}

// TransferEntry структура ответа на запрос к /coins/transfer
type TransferEntry struct {
	Balance int `json:"balance"`
}

// PurchaseEntry структура ответа на запрос к /merch/buy
type PurchaseEntry struct {
	Item    string `json:"item"`
	Balance int    `json:"balance"`
}
