package models

type Item struct {
	Name string `json:"name"`
	// MerchId int    `json:"merch_id"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}

type CoinsEntry struct {
	Date        string `json:"change_date"`
	CoinsBefore int    `json:"coins_before"`
	CoinsAfter  int    `json:"coins_after"`
}
