package service_test

import (
	"context"
	"testing"

	"merch_service/internal/models"
	"merch_service/internal/service"
	"merch_service/test/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserServiceRegister - проверяет метод Register у UserService :
// - корректное пробрасывание при первой регистрации
// - возвращение ошибки при наличии пользователя в бд
func TestUserServiceRegister(t *testing.T) {
	tests := []struct {
		name     string
		loginReq *models.LoginRequest
		wantErr  error
	}{
		{
			name: "успешная регистрация",
			loginReq: &models.LoginRequest{
				Login:    "newuser",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name: "пользователь уже существует",
			loginReq: &models.LoginRequest{
				Login:    "existing",
				Password: "password",
			},
			wantErr: models.ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userStorage := mock.NewMockUserStorage()
			purchaseStorage := mock.NewMockPurchaseStorage()
			coinsStorage := mock.NewMockCoinsStorage()
			userService := service.NewUserService(userStorage, purchaseStorage, coinsStorage)

			// Создаем существующего пользователя для второго теста
			if tt.wantErr == models.ErrUserExists {
				userStorage.Create(ctx, &models.User{
					Login:    "existing",
					Password: "password",
				})
			}

			err := userService.Register(ctx, tt.loginReq)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// TestUserServiceLogin - проверяет метод Login у UserService :
// - корректное пробрасывание при авторизации
// - возвращение ошибки при переданном неправильном пароле
// - возвращение ошибки при отсутствие пользователя в бд
func TestUserServiceLogin(t *testing.T) {
	tests := []struct {
		name     string
		loginReq *models.LoginRequest
		wantErr  error
	}{
		{
			name: "успешный вход",
			loginReq: &models.LoginRequest{
				Login:    "testuser",
				Password: "correct",
			},
			wantErr: nil,
		},
		{
			name: "неправильный пароль",
			loginReq: &models.LoginRequest{
				Login:    "testuser",
				Password: "wrong",
			},
			wantErr: models.ErrWrongPassword,
		},
		{
			name: "пользователь не найден",
			loginReq: &models.LoginRequest{
				Login:    "nonexistent",
				Password: "password",
			},
			wantErr: models.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userStorage := mock.NewMockUserStorage()
			purchaseStorage := mock.NewMockPurchaseStorage()
			coinsStorage := mock.NewMockCoinsStorage()
			userService := service.NewUserService(userStorage, purchaseStorage, coinsStorage)

			// Create test user
			userStorage.Create(ctx, &models.User{
				Login:    "testuser",
				Password: "correct",
			})

			err := userService.Login(ctx, tt.loginReq)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// TestMerchServiceHistory проверяет добавление в PurchaseHistory и в CoinsHistory после покупки
func TestMerchServiceHistory(t *testing.T) {
	ctx := context.Background()

	userStorage := mock.NewMockUserStorage()
	merchStorage := mock.NewMockMerchStorage()
	purchaseStorage := mock.NewMockPurchaseStorage()
	coinsStorage := mock.NewMockCoinsStorage()
	userService := service.NewUserService(userStorage, purchaseStorage, coinsStorage)
	merchService := service.NewMerchService(merchStorage, userStorage, purchaseStorage, coinsStorage)

	err := userStorage.Create(ctx, &models.User{
		Login:    "testuser",
		Password: "password",
		Coins:    1000,
	})
	assert.NoError(t, err)

	_, err = merchService.Buy(ctx, "testuser", "Футболка", 2)
	assert.NoError(t, err)

	coinsHistory, err := userService.CoinsHistory(ctx, "testuser")
	assert.NoError(t, err)
	assert.Len(t, coinsHistory, 1)
	assert.Equal(t, 1000, coinsHistory[0].CoinsBefore)
	assert.Equal(t, 800, coinsHistory[0].CoinsAfter)

	purchaseHistory, err := userService.PurchaseHistory(ctx, "testuser")
	assert.NoError(t, err)
	assert.Len(t, purchaseHistory, 1)
	assert.Equal(t, "Футболка", purchaseHistory[0].ItemName)
	assert.Equal(t, 2, purchaseHistory[0].Count)
}

// TestMerchServiceMerchList проверяет корректность возврата списка мерча
func TestMerchServiceMerchList(t *testing.T) {
	ctx := context.Background()
	merchStorage := mock.NewMockMerchStorage()
	userStorage := mock.NewMockUserStorage()

	purchaseStorage := mock.NewMockPurchaseStorage()
	coinsStorage := mock.NewMockCoinsStorage()
	merchService := service.NewMerchService(merchStorage, userStorage, purchaseStorage, coinsStorage)

	items, err := merchService.MerchList(ctx)
	// Не очень хороший тест, поскольку обход
	// мапы делается в случайном порядке
	// Поэтому от запуска к запуску результат меняется
	assert.NoError(t, err)
	assert.Len(t, items, 2)
	itemNames := make([]string, 0, len(items))
	for _, item := range items {
		itemNames = append(itemNames, item.Name)
	}
	assert.ElementsMatch(t, []string{"Кружка", "Футболка"}, itemNames)
}

// -------------------------- Проверка транзакций -----------------------------

// TestMerchServiceBuy - проверяет следующие сценарии работы метода Buy в MerchService:
// - успешная покупка
// - нехватка денег у пользователя для покупки
// - отсутсвие количества на складе
func TestMerchServiceBuy(t *testing.T) {
	tests := []struct {
		name       string
		userLogin  string
		merchName  string
		count      int
		wantCoins  int
		wantErr    error
		setupUser  func(s *mock.MockUserStorage)
		setupMerch func(s *mock.MockMerchStorage)
	}{
		{
			name:      "Успешная покупка",
			userLogin: "testuser",
			merchName: "Футболка",
			count:     2,
			wantCoins: 800,
			wantErr:   nil,
			setupUser: func(s *mock.MockUserStorage) {
				s.Create(context.Background(), &models.User{
					Login: "testuser",
					Coins: 1000,
				})
			},
			setupMerch: func(s *mock.MockMerchStorage) {},
		},
		{
			name:      "Не хватает деняг",
			userLogin: "testuser",
			merchName: "Футболка",
			count:     11,
			wantCoins: 1000,
			wantErr:   models.ErrNotEnoughCoins,
			setupUser: func(s *mock.MockUserStorage) {
				s.Create(context.Background(), &models.User{
					Login: "testuser",
					Coins: 1000,
				})
			},
			setupMerch: func(s *mock.MockMerchStorage) {},
		},
		{
			name:      "Не хватает мерча на складе",
			userLogin: "testuser",
			merchName: "Кружка",
			count:     10,
			wantCoins: 0,
			wantErr:   models.ErrNotEnoughMerch,
			setupUser: func(s *mock.MockUserStorage) {
				s.Create(context.Background(), &models.User{
					Login: "testuser",
					Coins: 1000,
				})
			},
			setupMerch: func(s *mock.MockMerchStorage) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userStorage := mock.NewMockUserStorage()
			merchStorage := mock.NewMockMerchStorage()
			purchaseStorage := mock.NewMockPurchaseStorage()
			coinsStorage := mock.NewMockCoinsStorage()
			merchService := service.NewMerchService(merchStorage, userStorage, purchaseStorage, coinsStorage)

			tt.setupUser(userStorage)
			tt.setupMerch(merchStorage)

			coins, err := merchService.Buy(ctx, tt.userLogin, tt.merchName, tt.count)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantCoins, coins)

				// Проверяем, что баланс действительно изменился
				user, err := userStorage.GetByLogin(ctx, tt.userLogin)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCoins, user.Coins)
			}
		})
	}
}

// TestTransactionServiceSend - проверяет метод Send в TransactionService на следующщие сценарии:
// - успешный перевод
// - отправителя нет в базе данных
// - недостаточно средств у отправителя
// - перевод самому себе
// - перевод отрицательной суммы
func TestTransactionServiceSend(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name              string
		senderLogin       string
		receiverLogin     string
		amount            int
		prepare           func(*mock.MockUserStorage)
		wantErr           error
		wantSenderCoins   int // -1 означает не проверять
		wantReceiverCoins int // -1 означает не проверять
	}{
		{
			name:          "Успешный перевод",
			senderLogin:   "sender",
			receiverLogin: "receiver",
			amount:        300,
			prepare: func(us *mock.MockUserStorage) {
				us.Create(ctx, &models.User{Login: "sender", Coins: 1000})
				us.Create(ctx, &models.User{Login: "receiver", Coins: 500})
			},
			wantErr:           nil,
			wantSenderCoins:   700,
			wantReceiverCoins: 800,
		},
		{
			name:          "Отправитель не найден",
			senderLogin:   "unknown",
			receiverLogin: "receiver",
			amount:        100,
			prepare: func(us *mock.MockUserStorage) {
				us.Create(ctx, &models.User{Login: "receiver", Coins: 500})
			},
			wantErr:           models.ErrUserNotFound,
			wantSenderCoins:   -1,
			wantReceiverCoins: 500,
		},
		{
			name:          "Недостаточно средств",
			senderLogin:   "sender",
			receiverLogin: "receiver",
			amount:        1500,
			prepare: func(us *mock.MockUserStorage) {
				us.Create(ctx, &models.User{Login: "sender", Coins: 1000})
				us.Create(ctx, &models.User{Login: "receiver", Coins: 500})
			},
			wantErr:           models.ErrNotEnoughCoins,
			wantSenderCoins:   1000,
			wantReceiverCoins: 500,
		},
		{
			name:          "Перевод самому себе",
			senderLogin:   "sender",
			receiverLogin: "sender",
			amount:        100,
			prepare: func(us *mock.MockUserStorage) {
				us.Create(ctx, &models.User{Login: "sender", Coins: 1000})
			},
			wantErr:           models.ErrSameSenderReceiver,
			wantSenderCoins:   1000,
			wantReceiverCoins: 1000,
		},
		{
			name:          "Отрицательная сумма",
			senderLogin:   "sender",
			receiverLogin: "receiver",
			amount:        -100,
			prepare: func(us *mock.MockUserStorage) {
				us.Create(ctx, &models.User{Login: "sender", Coins: 1000})
				us.Create(ctx, &models.User{Login: "receiver", Coins: 500})
			},
			wantErr:           models.ErrInvalidAmount,
			wantSenderCoins:   1000,
			wantReceiverCoins: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userStorage := mock.NewMockUserStorage()
			transactionStorage := mock.NewMockTransactionStorage()
			coinsStorage := mock.NewMockCoinsStorage()
			service := service.NewTransactionService(transactionStorage, userStorage, coinsStorage)

			if tc.prepare != nil {
				tc.prepare(userStorage)
			}

			err := service.Send(ctx, tc.senderLogin, tc.receiverLogin, tc.amount)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
			}

			if tc.wantSenderCoins > 0 {
				sender, err := userStorage.GetByLogin(ctx, tc.senderLogin)
				require.NoError(t, err)
				assert.Equal(t, tc.wantSenderCoins, sender.Coins)
			}

			if tc.wantReceiverCoins > 0 {
				receiver, err := userStorage.GetByLogin(ctx, tc.receiverLogin)
				require.NoError(t, err)
				assert.Equal(t, tc.wantReceiverCoins, receiver.Coins)
			}

			if tc.wantErr == nil {
				sender, err := userStorage.GetByLogin(ctx, tc.senderLogin)
				require.NoError(t, err)

				senderHistory, err := coinsStorage.Get(ctx, sender)
				require.NoError(t, err)
				require.Len(t, senderHistory, 1)
				assert.Equal(t, sender.Coins+tc.amount, senderHistory[0].CoinsBefore)
				assert.Equal(t, sender.Coins, senderHistory[0].CoinsAfter)

				receiver, err := userStorage.GetByLogin(ctx, tc.receiverLogin)
				require.NoError(t, err)

				receiverHistory, err := coinsStorage.Get(ctx, receiver)
				require.NoError(t, err)
				require.Len(t, receiverHistory, 1)
				assert.Equal(t, receiver.Coins-tc.amount, receiverHistory[0].CoinsBefore)
				assert.Equal(t, receiver.Coins, receiverHistory[0].CoinsAfter)
			}
		})
	}
}
