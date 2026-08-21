package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := LoadConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			if callback.Data == "catalog" {
				text := buildCatalogText()
				msg := tgbotapi.NewMessage(callback.Message.Chat.ID, text)
				bot.Send(msg)
			}
			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в PGT70!")
			msg.ReplyMarkup = mainMenuKeyboard()
			bot.Send(msg)
		}
	}
}
