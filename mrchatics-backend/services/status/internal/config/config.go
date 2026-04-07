package config

import (
	"os"
)

type Config struct {
	ServerPort      string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	CacheServiceURL string
}

func Load() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8081"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "status_service"),
		CacheServiceURL: getEnv("CACHE_SERVICE_URL", "http://localhost:8082"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
