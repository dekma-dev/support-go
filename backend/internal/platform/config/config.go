package config

import "os"

type Config struct {
	Environment               string
	HTTPPort                  string
	DatabaseURL               string
	KafkaBrokers              string
	NotificationConsumerGroup string
}

func Load() Config {
	return Config{
		Environment:               getEnv("APP_ENV", "local"),
		HTTPPort:                  getEnv("HTTP_PORT", "8080"),
		DatabaseURL:               getEnv("DATABASE_URL", ""),
		KafkaBrokers:              getEnv("KAFKA_BROKERS", ""),
		NotificationConsumerGroup: getEnv("KAFKA_NOTIFICATION_GROUP", "support-go-notification-worker"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
