package service

import (
	"context"
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/storage"
)

type TransactionServiceInterface interface {
	// Send - отправляет amount монет от sender к recv
	Send(ctx context.Context, sender, recv string, amount int) error
}

var _ TransactionServiceInterface = (*TransactionService)(nil)

// TransactionService - реализует интерфейс TransactionServiceInterface
type TransactionService struct {
	TransactionStorage storage.Storage
}

// Send - проверяет есть ли оба переданных пользователя,
// хватает ли денег отправителю для совершения операции,
// и совершает операцию отправки
func (t *TransactionService) Send(ctx context.Context, sender, recv string, amount int) error {
	sendUser, err := t.TransactionStorage.FindUserByLogin(ctx, sender)
	if err != nil {
		return err
	}

	recvUser, err := t.TransactionStorage.FindUserByLogin(ctx, recv)
	if err != nil {
		return err
	}

	if sendUser.Coins < amount {
		return models.ErrNotEnoughCoins
	}

	err = t.TransactionStorage.MakeTransaction(ctx, sendUser, recvUser, amount)
	if err != nil {
		return err
	}
	return nil
}
