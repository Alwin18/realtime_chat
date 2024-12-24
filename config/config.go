package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	SSLMode         string
	ServerPort      string
	SetMaxIdleConns string
	SetMaxOpenConns string
}

func LoadConfig() *Config {
	config := &Config{}
	envVars := map[string]*string{
		"DB_HOST":      &config.DBHost,
		"DB_PORT":      &config.DBPort,
		"DB_USER":      &config.DBUser,
		"DB_PASS":      &config.DBPass,
		"DB_NAME":      &config.DBName,
		"SERVER_PORT":  &config.ServerPort,
		"SET_MAX_IDLE": &config.SetMaxIdleConns,
		"SET_MAX_OPEN": &config.SetMaxOpenConns,
		"SSLMODE":      &config.SSLMode,
	}

	for key, ptr := range envVars {
		value := os.Getenv(key)
		if value == "" {
			log.Warnf("Missing environment variable: %s", key)
		}
		*ptr = value
	}

	return config
}
