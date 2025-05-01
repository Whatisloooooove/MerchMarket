package storagetest

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage"
	"merch_service/internal/storage/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
	ps "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestTransactionPG struct {
	suite.Suite
	pool               *pgxpool.Pool
	transactionStorage *postgres.TransactionPG
	ctx                context.Context
	container          testcontainers.Container
}

func (s *TestTransactionPG) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx

	pgContainer, err := ps.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		ps.WithDatabase("test_db"),
		ps.WithUsername("postgres"),
		ps.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(10*time.Second)),
	)
	require.NoError(s.T(), err)
	s.container = pgContainer

	connStr := "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable"

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(s.T(), err)
	s.pool = pool

	time.Sleep(2 * time.Second)

	_, err = pool.Exec(ctx, "SELECT 1")
	require.NoError(s.T(), err)

	dbconf := &storage.DBConfig{
		User:   "postgres",
		Pass:   "postgres",
		Addr:   "localhost",
		Port:   5432,
		DBName: "test_db",
	}

	err = storage.RunMigrations(dbconf, "../../migrations")
	require.NoError(s.T(), err)

	s.transactionStorage = postgres.NewTransactionStorage(pool)
}

func (s *TestTransactionPG) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}
	if s.container != nil {
		s.container.Terminate(context.Background())
	}
}

func (s *TestTransactionPG) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.transactions, merchshop.users, merchshop.coinhistory, merchshop.purchases CASCADE")
	require.NoError(s.T(), err)
}

func TestTestTransactionPG(t *testing.T) {
	suite.Run(t, new(TestTransactionPG))
}

// TestTransactionPGCreateTransaction - проверяет метод Create у TransactionPG:
// - успешная транзакция
// - нехватка деняг
// - невалидное значение
// - самому себе
// - отправитель не найден
// - получатель не найден
// - nil отправитель
// - nil получатель
func (s *TestTransactionPG) TestTransactionPGCreateTransaction() {
	t := s.T()

	sender := &models.User{
		Login:    "sender_user",
		Password: "sender_pass",
		Coins:    1000,
	}
	receiver := &models.User{
		Login:    "receiver_user",
		Password: "receiver_pass",
		Coins:    500,
	}

	userStorage := postgres.NewUserStorage(s.pool)
	err := userStorage.Create(s.ctx, sender)
	require.NoError(t, err)
	err = userStorage.Create(s.ctx, receiver)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		sender      *models.User
		receiver    *models.User
		amount      int
		wantErr     bool
		expectedErr error
		checkResult func(*testing.T, *models.User, *models.User)
	}{
		{
			name:     "успешная транзакция",
			sender:   sender,
			receiver: receiver,
			amount:   200,
			wantErr:  false,
			checkResult: func(t *testing.T, u *models.User, r *models.User) {
				updatedSender, err := userStorage.Get(s.ctx, u.Id)
				require.NoError(t, err)
				assert.Equal(t, 800, updatedSender.Coins)

				updatedReceiver, err := userStorage.Get(s.ctx, r.Id)
				require.NoError(t, err)
				assert.Equal(t, 700, updatedReceiver.Coins)

				var txCount int
				err = s.pool.QueryRow(s.ctx,
					"SELECT COUNT(*) FROM merchshop.transactions WHERE sender_id = $1 AND receiver_id = $2 AND amount = $3",
					u.Id, r.Id, 200).Scan(&txCount)
				require.NoError(t, err)
				assert.Equal(t, 1, txCount)
			},
		},
		{
			name:        "нехватка деняг",
			sender:      sender,
			receiver:    receiver,
			amount:      2000,
			wantErr:     true,
			expectedErr: models.ErrInsufficientCoins,
		},
		{
			name:        "невалидное значение",
			sender:      sender,
			receiver:    receiver,
			amount:      0,
			wantErr:     true,
			expectedErr: models.ErrInvalidAmount,
		},
		{
			name:        "самому себе",
			sender:      sender,
			receiver:    sender,
			amount:      100,
			wantErr:     true,
			expectedErr: models.ErrSameSenderReceiver,
		},
		{
			name:        "отправитель не найден",
			sender:      &models.User{Id: 9999},
			receiver:    receiver,
			amount:      100,
			wantErr:     true,
			expectedErr: models.ErrSenderNotFound,
		},
		{
			name:        "получатель не найден",
			sender:      sender,
			receiver:    &models.User{Id: 9999},
			amount:      100,
			wantErr:     true,
			expectedErr: models.ErrReceiverNotFound,
		},
		{
			name:        "nil отправитель",
			sender:      nil,
			receiver:    receiver,
			amount:      100,
			wantErr:     true,
			expectedErr: models.ErrEmptyTransaction,
		},
		{
			name:        "nil получатель",
			sender:      sender,
			receiver:    nil,
			amount:      100,
			wantErr:     true,
			expectedErr: models.ErrEmptyTransaction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.sender != nil && tc.sender.Id > 0 {
				_, err := s.pool.Exec(s.ctx,
					"UPDATE merchshop.users SET coins = $1 WHERE user_id = $2",
					1000, tc.sender.Id)
				require.NoError(t, err)
			}
			if tc.receiver != nil && tc.receiver.Id > 0 {
				_, err := s.pool.Exec(s.ctx,
					"UPDATE merchshop.users SET coins = $1 WHERE user_id = $2",
					500, tc.receiver.Id)
				require.NoError(t, err)
			}

			err := s.transactionStorage.Create(s.ctx, tc.sender, tc.receiver, tc.amount)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				if tc.checkResult != nil {
					tc.checkResult(t, tc.sender, tc.receiver)
				}
			}
		})
	}
}
