package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	TelegramBotToken string
	TelegramChatID string
	Port string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system env vars")
	}

	return &Config{
		DatabaseURL: mustGet("DATABASE_URL"),
		TelegramBotToken: mustGet("TELEGRAM_BOT_TOKEN"),
		TelegramChatID: mustGet("TELEGRAM_CHAT_ID"),
		Port: getOr("PORT", "8080"),
	}
}

func mustGet(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var: %s", key)
	}

	return v
}

func getOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}