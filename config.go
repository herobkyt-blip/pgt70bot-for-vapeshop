package main

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func LoadConfig() Config {
	godotenv.Load()
	return Config{
		BotToken: os.Getenv("BOT_TOKEN"),
	}
}
