package config

import "os"

type Config struct {
	Environment string
	HTTPPort    string
}

func Load() Config {
	return Config{
		Environment: getEnv("APP_ENV", "local"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

