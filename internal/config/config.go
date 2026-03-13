package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB   DB
	HTTP HTTP
	GRPC GRPC
}

type DB struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type HTTP struct {
	Port int
}

type GRPC struct {
	Port int
}

func Load() Config {
	return Config{
		DB: DB{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "reference_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		HTTP: HTTP{
			Port: getEnvInt("HTTP_PORT", 8082),
		},
		GRPC: GRPC{
			Port: getEnvInt("GRPC_PORT", 50051),
		},
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
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
