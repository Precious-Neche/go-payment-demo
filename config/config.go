package config

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	PaystackSecretKey   string
	PaystackPublicKey   string
	PaystackWebhookSecret string
	ServerPort          string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No .env file found, using system environment variables")
	}

	return &Config{
		PaystackSecretKey:   getEnv("PAYSTACK_SECRET_KEY", ""),
		PaystackPublicKey:   getEnv("PAYSTACK_PUBLIC_KEY", ""),
		PaystackWebhookSecret: getEnv("PAYSTACK_WEBHOOK_SECRET", ""),
		ServerPort:          getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}