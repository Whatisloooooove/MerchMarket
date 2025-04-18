package service

import (
	"context"
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type MerchServiceInterface interface {
	// Buy проверяет есть ли мерч в наличии и хватает
	// ли денег пользователю, после чего сохраняет данные
	Buy(ctx context.Context, login string, merchName string, count int) (int, error)

	MerchList(ctx context.Context) ([]models.Item, error)
}

type MerchService struct {
	MerchStorage storage.Storage
}
