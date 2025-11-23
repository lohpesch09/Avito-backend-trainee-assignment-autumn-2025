package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	LogLevel string
	DBURL    string
}

func NewConfig() *Config {
	godotenv.Load()
	return &Config{
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "debug"),
		DBURL:    getEnv("DBURL", ""),
	}
}

func getEnv(envName string, defaultValue string) string {
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}
	return defaultValue
}
