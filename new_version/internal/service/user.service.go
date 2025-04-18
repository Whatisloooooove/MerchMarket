package service

import (
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type UserServiceInterface interface {
	Login(logReq *models.LoginRequest) error

	Register(logReq *models.LoginRequest) error

	CoinsHistory(login string) ([]models.CoinsEntry, error)

	PurchaseHistory(login string) ([]models.PurchaseEntry, error)
}

type UserService struct {
	UserStorage storage.Storage
}
