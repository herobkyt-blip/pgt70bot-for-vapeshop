package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
			chatID := callback.Message.Chat.ID

			switch {
			case callback.Data == "catalog":
				for _, p := range products {
					text := fmt.Sprintf("%s - %.2f₽\n%s", p.Name, p.Price, p.Description)
					msg := tgbotapi.NewMessage(chatID, text)
					msg.ReplyMarkup = productKeyboard(p.ID)
					bot.Send(msg)
				}

			case callback.Data == "cart":
				cart := getCart(chatID)
				if len(cart) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Ваша корзина пуста."))
				} else {
					text := "Ваша корзина:\n"
					var total float64
					for _, p := range cart {
						text += fmt.Sprintf("%s - %.2f₽\n", p.Name, p.Price)
						total += p.Price
					}
					text += fmt.Sprintf("\nИтого: %.2f₽", total)
					bot.Send(tgbotapi.NewMessage(chatID, text))
				}
			case strings.HasPrefix(callback.Data, "add_"):
				idStr := strings.TrimPrefix(callback.Data, "add_")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					if product, found := findProduct(id); found {
						addToCart(chatID, product)
						bot.Send(tgbotapi.NewMessage(chatID, "Добавлено в корзину: "+product.Name))
					}
				}
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в PGT70")
			msg.ReplyMarkup = mainMenuKeyboard()
			bot.Send(msg)
		}
	}
}
