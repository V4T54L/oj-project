package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SERVER_PORT          int
	DB_URI               string
	JWT_SECRET           string
	TOKEN_EXPIRY_MINUTES int
}

// LoadEnv attempts to load .env file from the given path.
func LoadEnv(encPath string) error {
	if encPath == "" {
		encPath = ".env"
	}
	if err := godotenv.Load(encPath); err != nil {
		fmt.Println("Warning: No .env file found or failed to load:", err)
	}
	return nil
}

// NewConfig initializes config from environment variables with validation.
func NewConfig() (*Config, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080" // default fallback
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %w", err)
	}

	dbURI := os.Getenv("DB_DSN")
	if dbURI == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	tokenExpiryMinutesStr := os.Getenv("PORT")
	if tokenExpiryMinutesStr == "" {
		tokenExpiryMinutesStr = "8080" // default fallback
	}
	tokenExpiryMinutes, err := strconv.Atoi(tokenExpiryMinutesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid TOKEN_EXPIRY_MINUTES value: %w", err)
	}

	return &Config{
		SERVER_PORT:          port,
		DB_URI:               dbURI,
		JWT_SECRET:           jwtSecret,
		TOKEN_EXPIRY_MINUTES: tokenExpiryMinutes,
	}, nil
}
