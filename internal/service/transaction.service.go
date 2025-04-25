package service

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entites"
)

type TransactionServiceInterface interface {
	// Send - отправляет amount монет от sender к recv
	Send(ctx context.Context, sender, recv string, amount int) error
}

var _ TransactionServiceInterface = (*TransactionService)(nil)

// TransactionService - реализует интерфейс TransactionServiceInterface
type TransactionService struct {
	TransactionStorage entites.TransactionStorage
	UserStorage        entites.UserStorage
}

// NewTransactionService - создает объект TransactionService
func NewTransactionService(t entites.TransactionStorage, u entites.UserStorage) *TransactionService {
	return &TransactionService{
		TransactionStorage: t,
		UserStorage:        u,
	}
}

// Send - проверяет есть ли оба переданных пользователя,
// хватает ли денег отправителю для совершения операции,
// и совершает операцию отправки
func (t *TransactionService) Send(ctx context.Context, sender, recv string, amount int) error {
	sendUser, err := t.UserStorage.Get(ctx, sender)
	if err != nil {
		return err
	}

	recvUser, err := t.UserStorage.Get(ctx, recv)
	if err != nil {
		return err
	}

	if sendUser.Coins < amount {
		return models.ErrNotEnoughCoins
	}

	err = t.TransactionStorage.Create(ctx, sendUser, recvUser, amount)
	if err != nil {
		return err
	}
	return nil
}
