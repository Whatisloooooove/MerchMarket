// Файл отвечает за подключение к базе данных PostgreSQL с использованием GORM и автоматическую миграцию схем для моделей.
package db

import (
	"log"
	"merch_service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init инициализирует подключение к базе данных и выполняет миграцию схем.
// DSN: "host=localhost user=postgres password=postgres dbname=merch_store port=5432 sslmode=disable"
func Init() {
	dsn := "host=localhost user=postgres password=postgres dbname=merch_store port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Автоматическая миграция моделей: User и Wallet
	DB.AutoMigrate(&models.User{}, &models.Wallet{})
}
