package storagetest

import (
	"context"
	"fmt"
	"merch_service/internal/models"
	"merch_service/internal/storage"
	"merch_service/internal/storage/postgres"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestTransactionPG struct {
	suite.Suite
	pool               *pgxpool.Pool
	transactionStorage *postgres.TransactionPG
	ctx                context.Context
}

func (s *TestTransactionPG) SetupSuite() {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		"postgres", "postgres", "localhost", 5432, "postgres")

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		require.NoError(s.T(), err)
	}

	config.ConnConfig.TLSConfig = nil
	adminPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		require.NoError(s.T(), err)
	}

	_, err = adminPool.Exec(context.Background(),
		`
		DO
		$do$
		BEGIN
		IF EXISTS (
			SELECT FROM pg_user
			WHERE  usename = 'test_user') THEN

			RAISE NOTICE 'Role "my_user" already exists. Skipping.';
		ELSE
			CREATE USER test_user;
		END IF;
		END
		$do$;
		`)
	require.NoError(s.T(), err)

	testDBexists := false
	err = adminPool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'test_db');`).Scan(&testDBexists)
	require.NoError(s.T(), err)

	if !testDBexists {
		_, _ = adminPool.Exec(context.Background(),
			`CREATE DATABASE test_db OWNER test_user`)
	}

	_, err = adminPool.Exec(context.Background(),
		"GRANT ALL PRIVILEGES ON DATABASE test_db TO test_user")
	require.NoError(s.T(), err)

	adminPool.Close()

	dbconf := &storage.DBConfig{
		User:   "test_user",
		Pass:   "",
		Addr:   "localhost",
		Port:   5432,
		DBName: "test_db",
	}

	err = storage.CreateDb(dbconf)
	require.NoError(s.T(), err)

	connString = fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbconf.User,
		dbconf.Pass,
		dbconf.Addr,
		dbconf.Port,
		dbconf.DBName,
	)

	config, err = pgxpool.ParseConfig(connString)
	require.NoError(s.T(), err)
	config.ConnConfig.TLSConfig = nil

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	require.NoError(s.T(), err)

	err = storage.RunMigrations(dbconf, "../../migrations")
	require.NoError(s.T(), err)

	s.transactionStorage = postgres.NewTransactionStorage(pool)
	s.pool = pool
	s.ctx = context.Background()
}

func (s *TestTransactionPG) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()
	}

	adminPool, err := pgxpool.New(context.Background(),
		"user=postgres password= host=localhost port=5432 dbname=postgres sslmode=disable")
	if err != nil {
		s.T().Logf("Failed to connect as admin: %v", err)
		return
	}
	defer adminPool.Close()

	// Завершаем подключения к test_db
	_, err = adminPool.Exec(context.Background(), `
		SELECT pg_terminate_backend(pid)
		FROM pg_stat_activity
		WHERE datname = 'test_db' AND pid <> pg_backend_pid();
	`)
	if err != nil {
		s.T().Logf("Failed to terminate connections: %v", err)
	}

	// Удаляем базу и пользователя
	_, err = adminPool.Exec(context.Background(), "DROP DATABASE IF EXISTS test_db")
	if err != nil {
		s.T().Logf("Failed to drop test database: %v", err)
	}
}

func (s *TestTransactionPG) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.transactions CASCADE")
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
