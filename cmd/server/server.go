package main

import (
	"merch_service/internal/handlers"
	"merch_service/internal/server"
	"merch_service/internal/service"
	"merch_service/internal/storage"
	"merch_service/internal/storage/postgres"
)

func main() {
	db := storage.InitDB()

	// Инициализация базы данных (Storage interface) (можно заменить на свои моки)
	userStorage := postgres.NewUserStorage(db)
	merchStorage := postgres.NewMerchStorage(db)
	transactionStorage := postgres.NewTransactionStorage(db)
	coinsStorage := postgres.NewCoinsStorage(db)
	purchaseStorage := postgres.NewPurchaseStorage(db)

	// Инициализация сервисов
	userService := service.NewUserService(userStorage, purchaseStorage, coinsStorage)
	merchService := service.NewMerchService(merchStorage, userStorage, purchaseStorage, coinsStorage)
	transactionService := service.NewTransactionService(transactionStorage, userStorage, coinsStorage)

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userService)
	merchHandler := handlers.NewMerchHandler(merchService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Эти серивисы передаются в Server
	serv := server.NewMerchServer(userHandler, transactionHandler, merchHandler, "")
	serv.Start()
}
