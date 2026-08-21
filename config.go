package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
	AdminIDs map[int64]bool
}

func LoadConfig() Config {
	godotenv.Load()

	adminIDs := map[int64]bool{}
	idsRaw := os.Getenv("ADMIN_IDS")
	for _, part := range strings.Split(idsRaw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if id, err := strconv.ParseInt(part, 10, 64); err == nil {
			adminIDs[id] = true
		}
	}

	return Config{
		BotToken: os.Getenv("BOT_TOKEN"),
		AdminIDs: adminIDs,
	}
}

func (c Config) IsAdmin(userID int64) bool {
	return c.AdminIDs[userID]
}
