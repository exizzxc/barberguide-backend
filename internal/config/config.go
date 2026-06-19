package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
}

func Load() *Config {
	// Пробуем загрузить .env файл
	// Если его нет (прод) — просто читаем переменные окружения
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		Port:      port,
		DBDSN:     dbDSN,
		JWTSecret: jwtSecret,
	}
}
