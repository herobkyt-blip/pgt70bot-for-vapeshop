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
	loadCategories()

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
					msg := tgbotapi.NewMessage(chatID, text)
					msg.ReplyMarkup = cartKeyboard(cart)
					bot.Send(msg)
				}

			case callback.Data == "checkout":
				cart := getCart(chatID)
				if len(cart) > 0 {
					orderText := fmt.Sprintf("🆕 Новый заказ от @%s (ID: %d:\n", callback.From.UserName, chatID)
					var total float64
					for _, item := range cart {
						orderText += fmt.Sprintf("%s (%s) - %.2f₽\n", item.ProductName, item.Color, item.Price)
						total += item.Price
					}
					orderText += fmt.Sprintf("\nИтого: %.2f₽", total)

					for adminID := range cfg.AdminIDs {
						bot.Send(tgbotapi.NewMessage(adminID, orderText))
					}

					delete(carts, chatID)
					bot.Send(tgbotapi.NewMessage(chatID, "✅ Заказ оформлен! Скоро с вами свяжутся."))
				}

			case strings.HasPrefix(callback.Data, "removecart_"):
				indexStr := strings.TrimPrefix(callback.Data, "removecart_")
				index, err := strconv.Atoi(indexStr)
				if err == nil {
					cart := getCart(chatID)
					if index >= 0 && index < len(cart) {
						newCart := []CartItem{}
						for i, item := range cart {
							if i != index {
								newCart = append(newCart, item)
							}
						}
						carts[chatID] = newCart
					}
				}

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
					msg := tgbotapi.NewMessage(chatID, text)
					msg.ReplyMarkup = cartKeyboard(cart)
					bot.Send(msg)
				}

			case callback.Data == "back_main":
				msg := tgbotapi.NewMessage(chatID, "Добро пожаловать в BarbieBot")
				msg.ReplyMarkup = persistentKeyboard()
				bot.Send(msg)

			case callback.Data == "back_categories":
				msg := tgbotapi.NewMessage(chatID, "Выберите категорию:")
				msg.ReplyMarkup = categoriesKeyboard(categories)
				bot.Send(msg)

			case strings.HasPrefix(callback.Data, "back_products_"):
				category := strings.TrimPrefix(callback.Data, "back_products_")
				msg := tgbotapi.NewMessage(chatID, "Товары в категории "+category+":")
				msg.ReplyMarkup = productsInCategoryKeyboard(category)
				bot.Send(msg)

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
						if product.PhotoFileID != "" {
							photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(product.PhotoFileID))
							photo.Caption = text
							photo.ReplyMarkup = variantsKeyboard(product)
							bot.Send(photo)
						} else {
							msg := tgbotapi.NewMessage(chatID, text)
							msg.ReplyMarkup = variantsKeyboard(product)
							bot.Send(msg)
						}
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

			case callback.Data == "admin_add_category":
				drafts[chatID] = &ProductDraft{Step: StepWaitingCategoryName}
				bot.Send(tgbotapi.NewMessage(chatID, "Введите название категории:"))

			case strings.HasPrefix(callback.Data, "pickcat_"):
				if draft, exists := drafts[chatID]; exists && draft.Step == StepWaitingCategory {
					category := strings.TrimPrefix(callback.Data, "pickcat_")
					draft.Category = category
					draft.Step = StepWaitingColor
					bot.Send(tgbotapi.NewMessage(chatID, "Введите цвет (или напишите \"Стандарт\", если цвета нет):"))
				}

			case callback.Data == "confirm_variants":
				if draft, exists := drafts[chatID]; exists && draft.Step == StepWaitingColor {
					var variants []Variant
					for _, color := range draft.Colors {
						variants = append(variants, Variant{Color: color, Price: draft.Price, InStock: true})
					}

					newProduct := Product{
						ID:          len(products) + 1,
						Name:        draft.Name,
						Description: draft.Description,
						Category:    draft.Category,
						PhotoFileID: draft.PhotoFileID,
						Variants:    variants,
					}
					products = append(products, newProduct)
					saveProducts()

					delete(drafts, chatID)

					bot.Send(tgbotapi.NewMessage(chatID, "✅Товар добавлен: "+newProduct.Name))
				}

			case strings.HasPrefix(callback.Data, "delcat_"):
				categoryName := strings.TrimPrefix(callback.Data, "delcat_")

				newCategories := []string{}
				for _, c := range categories {
					if c != categoryName {
						newCategories = append(newCategories, c)
					}
				}
				categories = newCategories
				newProducts := []Product{}
				for _, p := range products {
					if p.Category != categoryName {
						newProducts = append(newProducts, p)
					}
				}
				products = newProducts

				saveCategories()
				saveProducts()

				bot.Send(tgbotapi.NewMessage(chatID, "🗑️Категория удалена: "+categoryName))

			case strings.HasPrefix(callback.Data, "delprod_"):
				idStr := strings.TrimPrefix(callback.Data, "delprod_")
				id, err := strconv.Atoi(idStr)
				if err == nil {
					newProducts := []Product{}
					for _, p := range products {
						if p.ID != id {
							newProducts = append(newProducts, p)
						}
					}
					products = newProducts
					saveProducts()
					bot.Send(tgbotapi.NewMessage(chatID, "🗑️Товар удалён."))
				}

			case callback.Data == "admin_main_menu":
				msg := tgbotapi.NewMessage(chatID, "Админ-панель")
				msg.ReplyMarkup = adminMenuKeyboard()
				bot.Send(msg)

			case callback.Data == "admin_categories_menu":
				msg := tgbotapi.NewMessage(chatID, "Управление категориями:")
				msg.ReplyMarkup = adminCategoriesMenuKeyboard()
				bot.Send(msg)

			case callback.Data == "admin_delete_category":
				if len(categories) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Нет категорий для удаления."))
				} else {
					msg := tgbotapi.NewMessage(chatID, "Выберите категорию для удаления:")
					msg.ReplyMarkup = deleteCategoryKeyboard()
					bot.Send(msg)
				}

			case callback.Data == "admin_products_menu":
				msg := tgbotapi.NewMessage(chatID, "Управление товарами:")
				msg.ReplyMarkup = adminProductsMenuKeyboard()
				bot.Send(msg)

			case callback.Data == "admin_delete_product":
				if len(products) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Нет товаров для удаления."))
				} else {
					msg := tgbotapi.NewMessage(chatID, "Выберите товар для удаления:")
					msg.ReplyMarkup = deleteProductKeyboad()
					bot.Send(msg)
				}
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать в BarbieBot")
			msg.ReplyMarkup = persistentKeyboard()
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

		if update.Message.Text == "📋 Каталог" {
			if len(categories) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Каталог пока пуст."))
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию:")
				msg.ReplyMarkup = categoriesKeyboard(categories)
				bot.Send(msg)
			}
		}

		if update.Message.Text == "🛒 Корзина" {
			cart := getCart(update.Message.Chat.ID)
			if len(cart) == 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша корзина пуста."))
			} else {
				text := "Ваша корзина:\n"
				var total float64
				for _, item := range cart {
					text += fmt.Sprintf("%s (%s) - %.2f₽\n", item.ProductName, item.Color, item.Price)
					total += item.Price
				}
				text += fmt.Sprintf("\nИтого: %.2f₽", total)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ReplyMarkup = cartKeyboard(cart)
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
		draft.Step = StepWaitingPhoto
		bot.Send(tgbotapi.NewMessage(chatID, "Отправьте фото товара:"))
	case StepWaitingPhoto:
		if len(message.Photo) == 0 {
			bot.Send(tgbotapi.NewMessage(chatID, "Это не похоже на фото. Отправьте фото товара:"))
			return
		}
		draft.PhotoFileID = message.Photo[len(message.Photo)-1].FileID
		draft.Step = StepWaitingCategory
		if len(categories) == 0 {
			bot.Send(tgbotapi.NewMessage(chatID, "Сначала добавьте хотя бы одну категорию через админ-панель."))
			delete(drafts, chatID)
			return
		}
		msg := tgbotapi.NewMessage(chatID, "Выберите категорию:")
		msg.ReplyMarkup = adminCategoriesKeyboard()
		bot.Send(msg)
	case StepWaitingColor:
		draft.Colors = append(draft.Colors, message.Text)

		text := "Добавлен цвет/вкус: " + strings.Join(draft.Colors, ", ") + ". Готово?"
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = confirmVariantsKeyboard()
		bot.Send(msg)

	case StepWaitingCategoryName:
		categories = append(categories, message.Text)
		saveCategories()
		delete(drafts, chatID)
		bot.Send(tgbotapi.NewMessage(chatID, "✅Категория добавлена: "+message.Text))
	}
}
