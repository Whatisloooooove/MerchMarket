package main

import "merch_service/new_version/internal/server"

func main() {

	// Инициализация базы данных (Storage interface) (можно заменить на свои моки)

	// Инициализированная база данных идет в userService, merchServie, transactionService
	// (можно заменить на свои моки)

	// Эти серивисы передаются в Server
	serv := server.NewMerchServer()
	serv.Start()
}
