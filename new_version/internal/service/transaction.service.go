package service

import (
	"context"
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

func (t *TransactionService) Send(ctx context.Context, sender, recv string, amount int) error {
	sendUser, err := t.TransactionStorage.FindUserByLogin(ctx, sender)
	if err != nil {
		return err
	}

	recvUser, err := t.TransactionStorage.FindUserByLogin(ctx, recv)
	if err != nil {
		return err
	}

	err = t.TransactionStorage.MakeTransaction(ctx, sendUser, recvUser, amount)
	if err != nil {
		return err
	}
	return nil
}
