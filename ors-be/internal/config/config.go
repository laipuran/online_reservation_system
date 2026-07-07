package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
	HTTPPort           string
	AllowedOrigins     string
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ors?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "dev-secret-do-not-use-in-production"),
		JWTExpirationHours: getEnvInt("JWT_EXPIRATION_HOURS", 24),
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		AllowedOrigins:     getEnv("ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
