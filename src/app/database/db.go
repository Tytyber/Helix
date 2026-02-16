package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func GetDBConfig() (string, error) {
	_ = godotenv.Load()

	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn, nil
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Проверяем обязательные поля
	if host == "" || user == "" || password == "" || dbname == "" {
		return "", fmt.Errorf("не заданы обязательные параметры подключения к БД (проверьте DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)")
	}

	// Если порт не указан, используем стандартный 5432
	if port == "" {
		port = "5432"
	}

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	return dsn, nil
}

func DBInit() (*sql.DB, error) {
	dsn, err := GetDBConfig()
	if err != nil {
		return nil, err
	}

	// Открываем соединение (драйвер зарегистрирован как "pgx")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Настраиваем пул соединений (опционально)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем связь с БД
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
