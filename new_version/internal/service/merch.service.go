package service

import (
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type MerchServiceInterface interface {
	// Buy проверяет есть ли мерч в наличии и хватает 
	// ли денег пользователю, после чего сохраняет данные
	Buy(login, merchName string) (int, error)

	MerchList() ([]models.Item, error)
}


type MerchService struct {
	MerchStorage storage.Storage

}