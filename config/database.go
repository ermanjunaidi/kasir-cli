package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// GetEnv mendapatkan environment variable dengan default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// InitDB menginisialisasi koneksi database
func InitDB() error {
	host := GetEnv("DB_HOST", "localhost")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "123123")
	dbname := GetEnv("DB_NAME", "kasir")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("gagal membuka koneksi database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("gagal koneksi ke database: %v", err)
	}

	return nil
}

// CloseDB menutup koneksi database
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
