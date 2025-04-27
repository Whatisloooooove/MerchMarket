package storage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Addr   string `yaml:"addr"`
	Port   int    `yaml:"port"`
	DBName string `yaml:"name"`
}

func loadConfig(configPath string) *DBConfig {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalln("не удалось определить абсолютный путь:", err.Error())
	}

	config, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalln("ошибка при чтении конфигурационного файла:", err.Error())
	}

	dc := DBConfig{}
	if err := yaml.Unmarshal(config, &dc); err != nil {
		log.Fatalln("ошибка при разборе yaml:", err.Error())
	}

	return &dc
}

func createDb(dbconf *DBConfig) error {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=postgres sslmode=disable",
		dbconf.User,
		dbconf.Pass,
		dbconf.Addr,
		dbconf.Port)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("ошибка парсинга конфига: %v", err)
	}
	config.ConnConfig.TLSConfig = nil
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к postgres: %v", err)
	}
	defer pool.Close()

	// Проверяем, существует ли БД
	var exists bool
	_ = pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		dbconf.DBName,
	).Scan(&exists)

	// Создаем БД если не существует
	if !exists {
		_, err = pool.Exec(context.Background(),
			fmt.Sprintf("CREATE DATABASE %s", dbconf.DBName),
		)
		if err != nil {
			return fmt.Errorf("не удалось создать БД: %v", err)
		}
		log.Printf("База данных '%s' создана", dbconf.DBName)
	}

	return nil
}

func runMigrations(dbconf *DBConfig) error {
	dbURL := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbconf.User, dbconf.Pass),
		Host:     fmt.Sprintf("%s:%d", dbconf.Addr, dbconf.Port),
		Path:     dbconf.DBName,
		RawQuery: "sslmode=disable",
	}

	migrationsAbsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalln("не удалось определить абсолютный путь миграций:", err.Error())
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsAbsPath),
		dbURL.String(),
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Миграции успешно применены!")
	return nil
}

func InitDB() *pgxpool.Pool {
	dbconf := loadConfig("configs/database_config.yml")
	log.Printf("DB Config: %+v", dbconf)
	if err := createDb(dbconf); err != nil {
		log.Fatalln("не удалось подключиться к базе данных:", err)
	}

	// в fmt.Sprintf добавить dbconf.Pass (по умолчанию если есть суперпользователь, пароль не нужно вводить)
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s",
		dbconf.User,
		dbconf.Pass,
		dbconf.Addr,
		dbconf.Port,
		dbconf.DBName)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}
	config.ConnConfig.TLSConfig = nil

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalln("не удалось подключиться к базе данных:", err)
	}

	if err := runMigrations(dbconf); err != nil {
		log.Fatalln("не удалось применить миграции:", err)
	}

	log.Printf("Успешное подключение к БД '%s'", dbconf.DBName)
	return pool
}
