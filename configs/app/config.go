package config

import (
	"fmt"
	"go-rest-api-boilerplate/types"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	PublicHost string
	Port       string
	Mode       string
	DBUser     string
	DBPassword string
	DBSource   string
	DBName     string
	Logger     types.LoggerSetting
}

const (
	twoDaysInSeconds = 60 * 60 * 24 * 2
)

var Envs = initConfig()

func initConfig() Config {
	logger := types.LoggerSetting{
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		FileLogName: getEnv("LOG_PATH", "text"),
		MaxBackups:  getEnvAsInt("LOG_MAX_BACKUPS", 3),
		MaxAge:      getEnvAsInt("LOG_MAX_AGE", 28),
		MaxSize:     getEnvAsInt("LOG_MAX_SIZE", 500),
		Compress:    getEnvAsBool("LOG_COMPRESS", true),
	}
	return Config{
		Mode:       getEnv("MODE", "dev"),
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "mypassword"),
		DBName:     getEnv("DB_NAME", "cars"),

		DBSource: getEnvOrError("DB_SOURCE"),
		Logger:   logger,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(value)
	}

	return fallback
}

func getEnvOrError(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("Environment variable %s is not set", key))

}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}

		return b
	}

	return fallback
}
