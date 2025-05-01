package storagetest

import (
	"context"
	"fmt"
	"merch_service/internal/models"
	"merch_service/internal/storage"
	"merch_service/internal/storage/postgres"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestMerchPG struct {
	suite.Suite
	pool         *pgxpool.Pool
	merchStorage *postgres.MerchPG
	ctx          context.Context
}

func (s *TestMerchPG) SetupSuite() {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost" // для CI
	}

	portStr := os.Getenv("POSTGRES_PORT")
	if portStr == "" {
		portStr = "5432" // для CI
	}

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		host,
		portStr,
		"postgres",
	)

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

	s.merchStorage = postgres.NewMerchStorage(pool)
	s.pool = pool
	s.ctx = context.Background()
}

func (s *TestMerchPG) TearDownSuite() {
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

func (s *TestMerchPG) SetupTest() {
	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.users, merchshop.coinhistory, merchshop.purchases CASCADE")
	require.NoError(s.T(), err)
}

func TestTestMerchPG(t *testing.T) {
	suite.Run(t, new(TestMerchPG))
}

// TestMerchPGCreate - тестирует функцию Create у MerchPG:
// - корректный мерч
// - пустое имя
// - некорректная стоимость
// - некорректное количество
// - нулевой мерч
func (s *TestMerchPG) TestMerchPGCreate() {
	t := s.T()

	testCases := []struct {
		name    string
		merch   *models.Item
		wantErr bool
		err     error
	}{
		{
			name: "корректный мерч",
			merch: &models.Item{
				Name:  "Test Item",
				Price: 1000,
				Stock: 10,
			},
			wantErr: false,
		},
		{
			name: "пустое имя",
			merch: &models.Item{
				Name:  "",
				Price: 1000,
				Stock: 10,
			},
			wantErr: true,
			err:     models.ErrEmptyMerchName,
		},
		{
			name: "некорректная стоимость",
			merch: &models.Item{
				Name:  "Test Item",
				Price: -100,
				Stock: 10,
			},
			wantErr: true,
			err:     models.ErrNegativePrice,
		},
		{
			name: "некорректное количество",
			merch: &models.Item{
				Name:  "Test Item",
				Price: 1000,
				Stock: -5,
			},
			wantErr: true,
			err:     models.ErrNegativeStock,
		},
		{
			name:    "нулевой мерч",
			merch:   nil,
			wantErr: true,
			err:     models.ErrEmptyMerch,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.merchStorage.Create(s.ctx, tc.merch)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.merch.Id)

				gotMerch, err := s.merchStorage.Get(s.ctx, tc.merch.Id)
				assert.NoError(t, err)
				assert.Equal(t, tc.merch.Name, gotMerch.Name)
				assert.Equal(t, tc.merch.Price, gotMerch.Price)
				assert.Equal(t, tc.merch.Stock, gotMerch.Stock)
			}
		})
	}
}

// TestMerchPGGet - проверяет Get у MerchPG:
// - успешное получение
// - получение несуществующего мерча
// - невалидный id
func (s *TestMerchPG) TestMerchPGGet() {
	t := s.T()

	testMerch := &models.Item{
		Name:  "Test Item",
		Price: 1000,
		Stock: 10,
	}
	err := s.merchStorage.Create(s.ctx, testMerch)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		id          int
		wantMerch   *models.Item
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "успешное получение",
			id:        testMerch.Id,
			wantMerch: testMerch,
			wantErr:   false,
		},
		{
			name:        "получение несуществующего мерча",
			id:          9999,
			wantMerch:   nil,
			wantErr:     true,
			expectedErr: models.ErrMerchNotFound,
		},
		{
			name:        "невалидный id",
			id:          0,
			wantMerch:   nil,
			wantErr:     true,
			expectedErr: models.ErrInvalidMerchID,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotMerch, err := s.merchStorage.Get(s.ctx, tc.id)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, gotMerch)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, gotMerch)
				assert.Equal(t, tc.wantMerch.Id, gotMerch.Id)
				assert.Equal(t, tc.wantMerch.Name, gotMerch.Name)
				assert.Equal(t, tc.wantMerch.Price, gotMerch.Price)
				assert.Equal(t, tc.wantMerch.Stock, gotMerch.Stock)
			}
		})
	}
}

// TesttMerchPGGetByName - тестирует метод GetByName у MerchPG:
// - успешное получение
// - получение несуществующего
func (s *TestMerchPG) TestMerchPGGetByName() {
	t := s.T()

	testMerch := &models.Item{
		Name:  "Unique Test Item",
		Price: 1000,
		Stock: 10,
	}
	err := s.merchStorage.Create(s.ctx, testMerch)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		merchName   string
		wantMerch   *models.Item
		wantErr     bool
		expectedErr error
	}{
		{
			name:      "успешное получение",
			merchName: "Unique Test Item",
			wantMerch: testMerch,
			wantErr:   false,
		},
		{
			name:        "получение несуществующего",
			merchName:   "Nonexistent Item",
			wantMerch:   nil,
			wantErr:     true,
			expectedErr: models.ErrMerchNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotMerch, err := s.merchStorage.GetByName(s.ctx, tc.merchName)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, gotMerch)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, gotMerch)
				assert.Equal(t, tc.wantMerch.Id, gotMerch.Id)
				assert.Equal(t, tc.wantMerch.Name, gotMerch.Name)
				assert.Equal(t, tc.wantMerch.Price, gotMerch.Price)
				assert.Equal(t, tc.wantMerch.Stock, gotMerch.Stock)
			}
		})
	}
}

// TestMerchPGUpdate - тестирует метод Update у MerchPG:
// - успешное обновление
// - пустое имя
// - некорректная стоимость
// - некорректное количество
// - нулевой мерч
func (s *TestMerchPG) TestMerchPGUpdate() {
	t := s.T()

	// Создаем тестовый товар
	originalMerch := &models.Item{
		Name:  "Original Item",
		Price: 1000,
		Stock: 10,
	}
	err := s.merchStorage.Create(s.ctx, originalMerch)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		modify      func(*models.Item) // Функция для изменения товара
		wantErr     bool
		expectedErr error
	}{
		{
			name: "успешное обновление",
			modify: func(m *models.Item) {
				m.Name = "Updated Item"
				m.Price = 2000
				m.Stock = 20
			},
			wantErr: false,
		},
		{
			name: "обновление с пустым именем",
			modify: func(m *models.Item) {
				m.Name = ""
			},
			wantErr:     true,
			expectedErr: models.ErrEmptyMerchName,
		},
		{
			name: "обновление с некорректной ценой",
			modify: func(m *models.Item) {
				m.Price = -100
			},
			wantErr:     true,
			expectedErr: models.ErrNegativePrice,
		},
		{
			name: "обновление с некорректным количеством",
			modify: func(m *models.Item) {
				m.Stock = -5
			},
			wantErr:     true,
			expectedErr: models.ErrNegativeStock,
		},
		{
			name: "обновление несуществующего мерча",
			modify: func(m *models.Item) {
				m.Id = 9999
			},
			wantErr:     true,
			expectedErr: models.ErrMerchNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			merchToUpdate := *originalMerch
			tc.modify(&merchToUpdate)

			err := s.merchStorage.Update(s.ctx, &merchToUpdate)

			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)

				updatedMerch, err := s.merchStorage.Get(s.ctx, originalMerch.Id)
				assert.NoError(t, err)
				assert.Equal(t, merchToUpdate.Name, updatedMerch.Name)
				assert.Equal(t, merchToUpdate.Price, updatedMerch.Price)
				assert.Equal(t, merchToUpdate.Stock, updatedMerch.Stock)
			}
		})
	}
}

// TestMerchPGGetList - тестирует метод GetList у MerchPG
func (s *TestMerchPG) TestMerchPGGetList() {
	t := s.T()

	_, err := s.pool.Exec(s.ctx, "TRUNCATE TABLE merchshop.merch CASCADE")
	require.NoError(t, err)

	items := []*models.Item{
		{Name: "Item 1", Price: 1000, Stock: 10},
		{Name: "Item 2", Price: 2000, Stock: 20},
		{Name: "Item 3", Price: 3000, Stock: 30},
	}

	for _, item := range items {
		err := s.merchStorage.Create(s.ctx, item)
		require.NoError(t, err)
	}

	gotItems, err := s.merchStorage.GetList(s.ctx)
	assert.NoError(t, err)
	assert.Len(t, gotItems, len(items))

	for _, expected := range items {
		found := false
		for _, actual := range gotItems {
			if expected.Id == actual.Id {
				assert.Equal(t, expected.Name, actual.Name)
				assert.Equal(t, expected.Price, actual.Price)
				assert.Equal(t, expected.Stock, actual.Stock)
				found = true
				break
			}
		}
		assert.True(t, found, "Item with ID %d not found in list", expected.Id)
	}
}
