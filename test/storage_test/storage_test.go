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

type TestUserPG struct {
	suite.Suite
	pool        *pgxpool.Pool
	usesStorage *postgres.UserPG
	ctx         context.Context
}

func (u *TestUserPG) SetupSuite() {
	connString := "user=postgres host=localhost port=5432 dbname=postgres sslmode=disable"

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		require.NoError(u.T(), err)
	}

	config.ConnConfig.TLSConfig = nil
	adminPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		require.NoError(u.T(), err)
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
			CREATE USER test_user LOGIN PASSWORD 'test_password';
		END IF;
		END
		$do$;
		`)
	require.NoError(u.T(), err)

	testDBexists := false
	err = adminPool.QueryRow(context.Background(),
		`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'test_db');`).Scan(&testDBexists)
	require.NoError(u.T(), err)

	if !testDBexists {
		_, _ = adminPool.Exec(context.Background(),
			`CREATE DATABASE test_db OWNER test_user`)
	}

	_, err = adminPool.Exec(context.Background(),
		"GRANT ALL PRIVILEGES ON DATABASE test_db TO test_user")
	require.NoError(u.T(), err)

	adminPool.Close()

	dbconf := &storage.DBConfig{
		User:   "test_user",
		Pass:   "test_password",
		Addr:   "localhost",
		Port:   5432,
		DBName: "test_db",
	}

	err = storage.CreateDb(dbconf)
	require.NoError(u.T(), err)

	connString = fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbconf.User,
		dbconf.Pass,
		dbconf.Addr,
		dbconf.Port,
		dbconf.DBName,
	)

	pool, err := pgxpool.New(u.ctx, connString)
	require.NoError(u.T(), err)
	u.pool = pool

	err = storage.RunMigrations(dbconf, "../../migrations")
	require.NoError(u.T(), err)

	u.usesStorage = postgres.NewUserStorage(pool)
	u.ctx = context.Background()
}

func (s *TestUserPG) TearDownSuite() {
	if s.pool != nil {
		s.pool.Close()

		adminPool, err := pgxpool.New(context.Background(),
			"user=postgres password=postgres host=localhost port=5432 dbname=postgres sslmode=disable")
		if err != nil {
			s.T().Logf("Failed to connect as admin: %v", err)
			return
		}
		defer adminPool.Close()

		_, err = adminPool.Exec(context.Background(),
			"DROP DATABASE IF EXISTS test_db")
		if err != nil {
			s.T().Logf("Failed to drop test database: %v", err)
		}

		_, err = adminPool.Exec(context.Background(),
			"DROP USER IF EXISTS test_user")
		if err != nil {
			s.T().Logf("Failed to drop test user: %v", err)
		}
	}
}

func (s *TestUserPG) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.users, merchshop.coinhistory, merchshop.purchases CASCADE")
	require.NoError(s.T(), err)
}

func TestTestUserPG(t *testing.T) {
	userPG := new(TestUserPG)
	userPG.ctx = context.Background()
	suite.Run(t, userPG)
}

// TestCreate - тестирует создание пользователя
func (u *TestUserPG) TestUserStorageCreate() {
	t := u.T()

	testCases := []struct {
		name    string
		user    *models.User
		wantErr bool
		err     error
	}{
		{
			name: "валидный юзер",
			user: &models.User{
				Login:    "testuser",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "пустой логин",
			user: &models.User{
				Login:    "",
				Password: "password123",
			},
			wantErr: true,
			err:     models.ErrEmptyUserLogin,
		},
		{
			name: "пустой пароль",
			user: &models.User{
				Login:    "testuser",
				Password: "",
			},
			wantErr: true,
			err:     models.ErrEmptyUserPassword,
		},
		{
			name:    "nil юзер",
			user:    nil,
			wantErr: true,
			err:     models.ErrEmptyUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := u.usesStorage.Create(u.ctx, tc.user)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.user.Id)
			}
		})
	}
}
