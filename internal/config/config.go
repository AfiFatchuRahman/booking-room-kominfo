package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	AppPort   string
	DBConnStr string
	JWTSecret string
}

func Load() Config {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "userposgres")
	pass := getEnv("DB_PASSWORD", "postgrespassword")
	name := getEnv("DB_NAME", "booking_room")
	ssl := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, pass, name, ssl,
	)

	return Config{
		AppPort:   getEnv("APP_PORT", "8080"),
		DBConnStr: dsn,
		JWTSecret: getEnv("JWT_SECRET", "diskominfo-diy-backend-secret-key-2026"),
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	for i := 0; i < 15; i++ {
		if err = db.Ping(); err == nil {
			return db, nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("gagal koneksi ke postgres: %w", err)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
