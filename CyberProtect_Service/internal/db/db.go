package db

import (
	"merch_service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := "host=localhost user=postgres password=postgres dbname=merch_store port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Run auto-migrations for your models
	DB.AutoMigrate(&models.User{}, &models.Wallet{})
}
