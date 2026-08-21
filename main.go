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
	loadProducts()

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
				categories := map[string]bool{}
				for _, p := range products {
					categories[p.Category] = true
				}
				if len(categories) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Каталог пока пуст."))
				} else {
					msg := tgbotapi.NewMessage(chatID, "Выберите категорию:")
					msg.ReplyMarkup = categoriesKeyboard(categories)
					bot.Send(msg)
				}

			case callback.Data == "cart":
				cart := getCart(chatID)
				if len(cart) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Ваша корзина пуста."))
				} else {
					text := "Ваша корзина:\n"
					var total float64
					for _, item := range cart {
						text += fmt.Sprintf("%s (%s) - %.2f₽\n", item.ProductName, item.Color, item.Price)
						total += item.Price
					}
					text += fmt.Sprintf("\nИтого: %.2f₽", total)
					bot.Send(tgbotapi.NewMessage(chatID, text))
				}

			case strings.HasPrefix(callback.Data, "category_"):
				category := strings.TrimPrefix(callback.Data, "category_")
				msg := tgbotapi.NewMessage(chatID, "Товары в категории "+category+":")
				msg.ReplyMarkup = productsInCategoryKeyboard(category)
				bot.Send(msg)

			case strings.HasPrefix(callback.Data, "product_"):
				idStr := strings.TrimPrefix(callback.Data, "product_")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					if product, found := findProduct(id); found {
						text := fmt.Sprintf("%s\n%s", product.Name, product.Description)
						msg := tgbotapi.NewMessage(chatID, text)
						msg.ReplyMarkup = variantsKeyboard(product)
						bot.Send(msg)
					}
				}

			case strings.HasPrefix(callback.Data, "variant_"):
				parts := strings.Split(strings.TrimPrefix(callback.Data, "variant_"), "_")
				if len(parts) == 2 {
					productID, err1 := strconv.Atoi(parts[0])
					variantIndex, err2 := strconv.Atoi(parts[1])
					if err1 == nil && err2 == nil {
						if product, found := findProduct(productID); found && variantIndex < len(product.Variants) {
							variant := product.Variants[variantIndex]
							addToCart(chatID, CartItem{ProductName: product.Name, Color: variant.Color, Price: variant.Price})
							bot.Send(tgbotapi.NewMessage(chatID, "Добавлено в корзину: "+product.Name+" ("+variant.Color+")"))
						}
					}
				}

			case callback.Data == "admin_add_product":
				drafts[chatID] = &ProductDraft{Step: StepWaitingName}
				bot.Send(tgbotapi.NewMessage(chatID, "Введите название товара:"))

			case callback.Data == "back_main":
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать в PGT70")
				msg.ReplyMarkup = mainMenuKeyboard()
				bot.Send(msg)

			case callback.Data == "back_categories":
				categories := map[string]bool{}
				for _, p := range products {
					categories[p.Category] = true
				}
				msg := tgbotapi.NewMessage(chatID, "Выберите категорию:")
				msg.ReplyMarkup = categoriesKeyboard(categories)
				bot.Send(msg)

			case strings.HasPrefix(callback.Data, "back_products_"):
				category := strings.TrimPrefix(callback.Data, "back_products_")
				msg := tgbotapi.NewMessage(chatID, "Товары в категории "+category+":")
				msg.ReplyMarkup = productsInCategoryKeyboard(category)
				bot.Send(msg)
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
		draft.Step = StepWaitingCategory
		bot.Send(tgbotapi.NewMessage(chatID, "Введите категорию:"))
	case StepWaitingCategory:
		draft.Category = message.Text
		draft.Step = StepWaitingColor
		bot.Send(tgbotapi.NewMessage(chatID, "Введите цвет (или напишите \"Стандарт\", если цвета нет):"))
	case StepWaitingColor:
		draft.Color = message.Text

		newProduct := Product{
			ID:          len(products) + 1,
			Name:        draft.Name,
			Description: draft.Description,
			Category:    draft.Category,
			Variants: []Variant{
				{Color: draft.Color, Price: draft.Price, InStock: true},
			},
		}
		products = append(products, newProduct)
		saveProducts()

		delete(drafts, chatID)

		bot.Send(tgbotapi.NewMessage(chatID, "✅Товар добавлен: "+newProduct.Name))
	}
}
