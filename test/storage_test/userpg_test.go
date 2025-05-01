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
	userStorage *postgres.UserPG
	ctx         context.Context
}

func (s *TestUserPG) SetupSuite() {
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

	s.userStorage = postgres.NewUserStorage(pool)
	s.pool = pool
	s.ctx = context.Background()
}

func (s *TestUserPG) TearDownSuite() {
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

func (s *TestUserPG) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.users, merchshop.coinhistory, merchshop.purchases CASCADE")
	require.NoError(s.T(), err)
}

func TestTestUserPG(t *testing.T) {
	userPG := new(TestUserPG)
	userPG.ctx = context.Background()
	suite.Run(t, userPG)
}

// TestCreate - тестирует создание пользователя:
// - валидный юзер
// - пустой юзер
// - пустой логин
// - пустой пароль
func (s *TestUserPG) TestUserStorageCreate() {
	t := s.T()

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
			err := s.userStorage.Create(s.ctx, tc.user)

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

// TestSerPGGet тестирует получение пользователя по ID:
// - успешное получение
// - получение несуществующего пользователя
// - получение невалидного пользователя
func (s *TestUserPG) TestGet() {
	t := s.T()

	user := &models.User{
		Login:    "testuser",
		Password: "password123",
		Coins:    1000,
	}
	err := s.userStorage.Create(s.ctx, user)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		userID      int
		wantUser    *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "успешное получение",
			userID:   user.Id,
			wantUser: user,
			wantErr:  false,
		},
		{
			name:        "получение несуществующего пользователя",
			userID:      9999,
			wantUser:    nil,
			wantErr:     true,
			expectedErr: models.ErrUserNotFound,
		},
		{
			name:        "получение невалидного пользователя",
			userID:      0,
			wantUser:    nil,
			wantErr:     true,
			expectedErr: models.ErrInvalidUserID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotUser, err := s.userStorage.Get(s.ctx, tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, gotUser)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, gotUser)
				assert.Equal(t, tc.wantUser.Id, gotUser.Id)
				assert.Equal(t, tc.wantUser.Login, gotUser.Login)
				assert.Equal(t, tc.wantUser.Password, gotUser.Password)
				assert.Equal(t, tc.wantUser.Coins, gotUser.Coins)
			}
		})
	}
}

// TestUpdate тестирует обновление пользователя
// - успешное обновление
// - обновление несуществующего
// - обновление с неверный id
func (s *TestUserPG) TestUserPGUpdate() {
	t := s.T()

	originalUser := &models.User{
		Login:    "testuser",
		Password: "password123",
		Coins:    100,
	}
	err := s.userStorage.Create(s.ctx, originalUser)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		prepare     func(*models.User) *models.User
		wantErr     bool
		expectedErr error
		checkResult func(*testing.T, *models.User)
	}{
		{
			name: "успешное обновление",
			prepare: func(u *models.User) *models.User {
				return &models.User{
					Id:       u.Id,
					Login:    "updateduser",
					Password: "newpassword",
					Coins:    200,
				}
			},
			wantErr: false,
			checkResult: func(t *testing.T, updated *models.User) {
				gotUser, err := s.userStorage.Get(s.ctx, updated.Id)
				require.NoError(t, err)
				assert.Equal(t, updated.Login, gotUser.Login)
				assert.Equal(t, updated.Password, gotUser.Password)
				assert.Equal(t, updated.Coins, gotUser.Coins)
			},
		},
		{
			name: "обновление несуществующего",
			prepare: func(u *models.User) *models.User {
				return &models.User{
					Id:       9999,
					Login:    "nonexistent",
					Password: "pass",
					Coins:    0,
				}
			},
			wantErr:     true,
			expectedErr: models.ErrUserNotFound,
			checkResult: func(t *testing.T, _ *models.User) {
			},
		},
		{
			name: "обновление с неверный id",
			prepare: func(u *models.User) *models.User {
				return &models.User{
					Id:       0,
					Login:    "invalid",
					Password: "pass",
					Coins:    0,
				}
			},
			wantErr:     true,
			expectedErr: models.ErrInvalidUserID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userToUpdate := tc.prepare(originalUser)

			err := s.userStorage.Update(s.ctx, userToUpdate)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if tc.checkResult != nil {
				tc.checkResult(t, userToUpdate)
			}
		})
	}
}

// TestUserPGGetByLogin тестирует получение пользователя по логину
// - успешное выполнение
// - получение несуществующего пользователя
// - пустой логие
func (s *TestUserPG) TestUserPGGetByLogin() {
	t := s.T()

	testUser := &models.User{
		Login:    "testuser",
		Password: "password123",
		Coins:    1000,
	}
	err := s.userStorage.Create(s.ctx, testUser)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		login       string
		wantUser    *models.User
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "успешное выполнение",
			login:    testUser.Login,
			wantUser: testUser,
			wantErr:  false,
		},
		{
			name:        "получение несуществующего пользователя",
			login:       "nonexistent",
			wantUser:    nil,
			wantErr:     true,
			expectedErr: models.ErrUserNotFound,
		},
		{
			name:        "пустой логин",
			login:       "",
			wantUser:    nil,
			wantErr:     true,
			expectedErr: models.ErrEmptyUserLogin,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotUser, err := s.userStorage.GetByLogin(s.ctx, tc.login)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, gotUser)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, gotUser)
				assert.Equal(t, tc.wantUser.Id, gotUser.Id)
				assert.Equal(t, tc.wantUser.Login, gotUser.Login)
				assert.Equal(t, tc.wantUser.Password, gotUser.Password)
				assert.Equal(t, tc.wantUser.Coins, gotUser.Coins)
			}
		})
	}
}
