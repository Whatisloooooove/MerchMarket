package main

import (
	"merch_service/internal/handlers"
	"merch_service/internal/server"
	"merch_service/internal/service"
	"merch_service/test/mock"
)

func main() {

	// Инициализация базы данных (Storage interface) (можно заменить на свои моки)
	userStorage := mock.NewMockUserStorage()
	merchStorage := mock.NewMockMerchStorage()
	transactionStorage := mock.NewMockTransactionStorage()

	// Инициализация сервисов
	userService := service.NewUserService(userStorage)
	merchService := service.NewMerchService(merchStorage)
	transactionService := service.NewTransactionService(transactionStorage, userStorage)

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userService)
	merchHandler := handlers.NewMerchHandler(merchService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Эти серивисы передаются в Server
	serv := server.NewMerchServer(userHandler, transactionHandler, merchHandler)
	serv.Start()
}
