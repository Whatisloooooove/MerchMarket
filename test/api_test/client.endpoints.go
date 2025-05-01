package apitest

import "merch_service/internal/models"

// ВАЖНО !!!
// TODO: добавить в README!
// API нашего проекта будет поддерживать следующий JSON формат:
// 		error_code  - код ошибки
// 		message - сообщение об ошибке (успехе)
// 		data - данные в формате задаваемом API
//
// Например, успешный запрос к /merch будет выглядеть так:
// {
//		error_code: 200,
//		message: "", # пустой либо "Merch list loaded successfully"
//		data: [
//			{
// 			name: "shoes",
// 			price: 100,
// 			stock: 20
// 			},
// 			...
//		]
//}
//
// [source](https://medium.com/@bojanmajed/standard-json-api-response-format-c6c1aabcaa6d)

// ResponseBody структура для хранения ответа сервера
// см. контракт описанный выше
type ResponseBody struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
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
	Balance int `json:"balance"`
}
