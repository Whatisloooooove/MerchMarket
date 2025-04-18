package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

type DB struct {
	pool   *pgxpool.Pool
	config *DBConfig
}

// Придется делать так, будет синглтоном
var db *DB

func loadConfig(configPath string) *DBConfig {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("не удалось получить путь до рабочей директорий:", err.Error())
	}

	// Поднимаемся в корень директорий и оттуда смотрим на configPath
	// TODO выглядит криво, что-то придумать!
	configPath = filepath.Join(filepath.Dir(filepath.Dir(wd)), configPath)
	config, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatalln("ошибка при чтений конфигурационного файла:", err.Error())
	}

	dbconf := DBConfig{}

	yaml.Unmarshal(config, &dbconf)
	return &dbconf
}

func InitDB() {
	dbconf := loadConfig("configs/database_config.yml")

	// в fmt.Sprintf добавить dbconf.Pass (по умолчанию если есть суперпользователь, пароль не нужно вводить)
	connString := fmt.Sprintf("postgres://%s@%s:%d/%s",
		dbconf.User,
		dbconf.Addr,
		dbconf.Port,
		dbconf.DBName)
	pool, err := pgxpool.New(context.TODO(), connString)

	if err != nil {
		log.Fatalln("не удалось подключиться к базе данных:", err)
	}

	db = &DB{
		pool:   pool,
		config: dbconf,
	}

	// var greeting string
	// err = db.db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(greeting)
}

// Connect - обертка для доступа к базе данных из вне
func Connect() *DB {
	return db
}

// Close - закрывает соединения в пуле
func (db *DB) Close() {
	db.pool.Close()
}
