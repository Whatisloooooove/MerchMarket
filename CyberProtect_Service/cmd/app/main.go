// Главная точка входа в приложение. Здесь происходит:
// - Инициализация подключения к базе данных.
// - Создание и настройка маршрутизатора (с использованием Gin).
// - Определение публичных и защищённых маршрутов.
package main

import (
	"log"
	"merch_service/internal/db"
	"merch_service/internal/handlers"
	"merch_service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация подключения к базе данных и миграция схем
	db.Init()

	// Создание нового маршрутизатора Gin
	r := gin.Default()

	// Публичные маршруты
	api := r.Group("/api")
	{
		// Маршрут для авторизации (логин)
		api.POST("/login", handlers.LoginHandler)
	}

	// Защищённые маршруты (требуют JWT аутентификации)
	auth := api.Group("/")
	auth.Use(middlewares.JWTAuth())
	{
		// Маршрут для получения информации о пользователе
		auth.GET("/me", handlers.MeHandler)
	}

	// Запуск HTTP-сервера на порту 8080
	log.Fatal(r.Run(":8080"))
}
