package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment               string
	HTTPPort                  string
	DatabaseURL               string
	JWTSecret                 string
	CORSOrigins               string
	KafkaBrokers              string
	NotificationConsumerGroup string
	NotificationRetryMax      int
	NotificationRetryBackoff  int
	NotificationDLQTopic      string
}

func Load() Config {
	return Config{
		Environment:               getEnv("APP_ENV", "local"),
		HTTPPort:                  getEnv("HTTP_PORT", "8080"),
		DatabaseURL:               getEnv("DATABASE_URL", ""),
		JWTSecret:                 getEnv("JWT_SECRET", ""),
		CORSOrigins:               getEnv("CORS_ORIGINS", "*"),
		KafkaBrokers:              getEnv("KAFKA_BROKERS", ""),
		NotificationConsumerGroup: getEnv("KAFKA_NOTIFICATION_GROUP", "support-go-notification-worker"),
		NotificationRetryMax:      getEnvInt("NOTIFICATION_RETRY_MAX", 3),
		NotificationRetryBackoff:  getEnvInt("NOTIFICATION_RETRY_BACKOFF_MS", 500),
		NotificationDLQTopic:      getEnv("KAFKA_NOTIFICATION_DLQ_TOPIC", "support.notification.dlq"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}

	return value
}
