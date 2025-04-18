package service

import "merch_service/new_version/internal/storage"

type TransactionServiceInterface interface {
	Send(sender, recv string, amount int) error
}

type TransactionService struct {
	TransactionStorage storage.Storage
}
