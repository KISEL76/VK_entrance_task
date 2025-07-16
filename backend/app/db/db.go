package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func newPostgres() (*sql.DB, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения с БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка соединения с БД: %w", err)
	}
	return db, nil
}
