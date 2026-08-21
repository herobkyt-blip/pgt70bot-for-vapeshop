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
			case callback.Data == "admin_add_product":
				drafts[chatID] = &ProductDraft{Step: StepWaitingName}
				bot.Send(tgbotapi.NewMessage(chatID, "Введите название товара:"))
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			continue
		}

		if update.Message == nil {
			continue
		}

		if draft, exists := drafts[update.Message.Chat.ID]; exists {
			handleAdminStep(bot, update.Message, draft)
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в PGT70")
			msg.ReplyMarkup = mainMenuKeyboard()
			bot.Send(msg)
		}

		if update.Message.Text == "/admin" {
			userID := update.Message.From.ID
			if !cfg.IsAdmin(userID) {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Доступ запрещён."))
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Админ-панель")
				msg.ReplyMarkup = adminMenuKeyboard()
				bot.Send(msg)
			}
		}
	}
}

func handleAdminStep(bot *tgbotapi.BotAPI, message *tgbotapi.Message, draft *ProductDraft) {
	chatID := message.Chat.ID

	switch draft.Step {
	case StepWaitingName:
		draft.Name = message.Text
		draft.Step = StepWaitingPrice
		bot.Send(tgbotapi.NewMessage(chatID, "Введите цену товара (числом):"))
	case StepWaitingPrice:
		price, err := strconv.ParseFloat(message.Text, 64)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(chatID, "Это не похоже на число. Введите цену ещё раз:"))
			return
		}
		draft.Price = price
		draft.Step = StepWaitingDesc
		bot.Send(tgbotapi.NewMessage(chatID, "Введите описание товара:"))
	case StepWaitingDesc:
		draft.Description = message.Text

		newProduct := Product{
			ID:          len(products) + 1,
			Name:        draft.Name,
			Price:       draft.Price,
			Description: draft.Description,
			InStock:     true,
		}
		products = append(products, newProduct)

		delete(drafts, chatID)

		bot.Send(tgbotapi.NewMessage(chatID, "✅Товар добавлен: "+newProduct.Name))
	}
}
