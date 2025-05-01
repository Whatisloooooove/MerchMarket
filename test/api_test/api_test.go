package apitest

import (
	"context"
	"log"
	"merch_service/internal/handlers"
	"merch_service/internal/models"
	"merch_service/internal/server"
	"merch_service/internal/service"
	"merch_service/test/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ServerStart(t *testing.T) server.Server {
	t.Helper()

	userStorage := mock.NewMockUserStorage()
	purchaseStorage := mock.NewMockPurchaseStorage()
	coinsStorage := mock.NewMockCoinsStorage()
	merchStorage := mock.NewMockMerchStorage()
	transactionStorage := mock.NewMockTransactionStorage()

	userService := service.NewUserService(userStorage, purchaseStorage, coinsStorage)
	merchService := service.NewMerchService(merchStorage, userStorage, purchaseStorage, coinsStorage)
	transactionService := service.NewTransactionService(transactionStorage, userStorage, coinsStorage)

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userService)
	merchHandler := handlers.NewMerchHandler(merchService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Эти серивисы передаются в Server
	// Захардкоженые пути, простите =(
	serv := server.NewMerchServer(userHandler, transactionHandler, merchHandler, "../../configs/server_config.yml")

	go serv.Start()

	return serv
}

func TestRegisterAPI(t *testing.T) {
	server := ServerStart(t)
	cli := NewClient()

	tests := []struct {
		name     string
		loginReq *models.LoginRequest
		response ResponseBody
	}{
		{
			name: "Регистрация нового пользователя",
			loginReq: &models.LoginRequest{
				Login:    "aboba",
				Password: "123123",
			},
			response: ResponseBody{
				ErrorCode: http.StatusOK,
				Message:   handlers.RegistrationOK,
			},
		},
		{
			name: "Попытка регистрации существующего пользователя",
			loginReq: &models.LoginRequest{
				Login:    "aboba",
				Password: "123123",
			},
			response: ResponseBody{
				ErrorCode: http.StatusBadRequest,
				Message:   handlers.UserExistsError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			response, err := cli.Register(context.Background(), tt.loginReq)

			require.NoError(t, err)
			assert.Equal(t, tt.response.ErrorCode, response.ErrorCode)
			assert.Equal(t, tt.response.Message, response.Message)
		})
	}

	server.Stop()
}

func TestLoginAPI(t *testing.T) {
	server := ServerStart(t)
	cli := NewClient()

	tests := []struct {
		name          string
		loginReq      *models.LoginRequest
		response      ResponseBody
		shouldSucceed bool
	}{
		{
			name: "Пользователя нет в базе",
			loginReq: &models.LoginRequest{
				Login:    "aboba",
				Password: "123123",
			},
			response: ResponseBody{
				ErrorCode: http.StatusBadRequest,
				Message:   handlers.UserNotFoundError,
			},
			shouldSucceed: false,
		},
		{
			name: "Успешная аутентификация",
			loginReq: &models.LoginRequest{
				Login:    "aboba",
				Password: "123123",
			},
			response: ResponseBody{
				ErrorCode: http.StatusOK,
				Message:   handlers.TokensOK,
			},
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldSucceed {
				// Имитация предварительной регистрации
				registerResp, err := cli.Register(context.Background(), tt.loginReq)

				require.Equal(t, http.StatusOK, registerResp.ErrorCode)
				require.Equal(t, handlers.RegistrationOK, registerResp.Message)
				require.NoError(t, err)
			}

			response, err := cli.GetTokens(context.Background(), tt.loginReq)

			require.NoError(t, err)
			assert.Equal(t, tt.response.ErrorCode, response.ErrorCode)
			assert.Equal(t, tt.response.Message, response.Message)

			// Проверяем что получили токены
			if tt.shouldSucceed {
				tokens, ok := response.Data.(*UserTokens)
				if !ok {
					log.Fatal("должны получить токены, получили не то")
				}

				assert.NotEmpty(t, tokens.Token)
				assert.NotEmpty(t, tokens.Refresh)
			}
		})
	}

	server.Stop()
}

func TestBuyAPI(t *testing.T) {
	server := ServerStart(t)
	cli := NewClient()

	tests := []struct {
		name          string
		loginReq      *models.LoginRequest
		purchReq      *models.PurchaseRequest
		response      ResponseBody
		shouldSucceed bool
	}{
		{
			name: "Обычная покупка",
			loginReq: &models.LoginRequest{
				Login:    "aboba",
				Password: "123123",
			},
			purchReq: &models.PurchaseRequest{
				Item:  "Футболка", // Взято из mock.NewMockMerchStorage
				Count: 1,
			},
			response: ResponseBody{
				ErrorCode: http.StatusOK,
				Message:   handlers.PurchaseOK,
				Data:      900, // Ожидаемый баланс после покупки (посчитано вручную, см. mock.MewMockMerchStorage)
			},
			shouldSucceed: true,
		},
		{
			name: "Недостаточно средств",
			loginReq: &models.LoginRequest{
				Login:    "aboba2",
				Password: "123123",
			},
			purchReq: &models.PurchaseRequest{
				Item:  "ОченьДорогаяВещь", // Добавлено в mock.NewMockMerchStorage
				Count: 1,
			},
			response: ResponseBody{
				ErrorCode: http.StatusBadRequest,
				Message:   handlers.NotEnoughCoinsError,
			},
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Для начала регистрируем пользователя
			_, err := cli.Register(context.Background(), tt.loginReq)
			require.NoError(t, err)

			response, err := cli.GetTokens(context.Background(), tt.loginReq)
			require.NoError(t, err)

			tokens, ok := response.Data.(*UserTokens)
			if !ok {
				log.Fatalln("Что то пошло не так, должны получить токены")
			}

			purchResp, err := cli.Buy(context.Background(), tt.purchReq, tokens)
			require.NoError(t, err)

			assert.Equal(t, tt.response.ErrorCode, purchResp.ErrorCode)
			assert.Equal(t, tt.response.Message, purchResp.Message)

			if tt.shouldSucceed {
				balance, ok := purchResp.Data.(*PurchaseEntry)
				if !ok {
					log.Fatalln("Что то пошло не так, должны получить число")
				}

				assert.Equal(t, tt.response.Data, balance.Balance)
			}
		})
	}

	server.Stop()
}

// func TestTransferAPI(t *testing.T) {
// 	server := ServerStart(t)
// 	cli := NewClient()

// 	tests := []struct {
// 		name          string
// 		user1         *models.LoginRequest
// 		user2         *models.LoginRequest
// 		purchReq      *models.PurchaseRequest
// 		response      ResponseBody
// 		shouldSucceed bool
// 	}{
// 		{
// 			name: "Обычная покупка",
// 			user1: &models.LoginRequest{
// 				Login:    "aboba",
// 				Password: "123123",
// 			},
// 			purchReq: &models.PurchaseRequest{
// 				Item:  "Футболка", // Взято из mock.NewMockMerchStorage
// 				Count: 1,
// 			},
// 			response: ResponseBody{
// 				ErrorCode: http.StatusOK,
// 				Message:   handlers.PurchaseOK,
// 				Data:      900, // Ожидаемый баланс после покупки (посчитано вручную, см. mock.MewMockMerchStorage)
// 			},
// 			shouldSucceed: true,
// 		},
// 		{
// 			name: "Недостаточно средств",
// 			loginReq: &models.LoginRequest{
// 				Login:    "aboba2",
// 				Password: "123123",
// 			},
// 			purchReq: &models.PurchaseRequest{
// 				Item:  "ОченьДорогаяВещь", // Добавлено в mock.NewMockMerchStorage
// 				Count: 1,
// 			},
// 			response: ResponseBody{
// 				ErrorCode: http.StatusBadRequest,
// 				Message:   handlers.NotEnoughCoinsError,
// 			},
// 			shouldSucceed: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Для начала регистрируем пользователя
// 			_, err := cli.Register(context.Background(), tt.loginReq)
// 			require.NoError(t, err)

// 			response, err := cli.GetTokens(context.Background(), tt.loginReq)
// 			require.NoError(t, err)

// 			tokens, ok := response.Data.(*UserTokens)
// 			if !ok {
// 				log.Fatalln("Что то пошло не так, должны получить токены")
// 			}

// 			purchResp, err := cli.Buy(context.Background(), tt.purchReq, tokens)
// 			require.NoError(t, err)

// 			assert.Equal(t, tt.response.ErrorCode, purchResp.ErrorCode)
// 			assert.Equal(t, tt.response.Message, purchResp.Message)

// 			if tt.shouldSucceed {
// 				balance, ok := purchResp.Data.(*PurchaseEntry)
// 				if !ok {
// 					log.Fatalln("Что то пошло не так, должны получить число")
// 				}

// 				assert.Equal(t, tt.response.Data, balance.Balance)
// 			}
// 		})
// 	}

// 	server.Stop()
// }
