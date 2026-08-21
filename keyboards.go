package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mainMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Каталог", "catalog"),
			tgbotapi.NewInlineKeyboardButtonData("Корзина", "cart"),
		),
	)
}

func adminMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить товар", "admin_add_product"),
		),
	)
}

func categoriesKeyboard(categories map[string]bool) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(category, "category_"+category)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func productsInCategoryKeyboard(category string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, p := range products {
		if p.Category == category {
			button := tgbotapi.NewInlineKeyboardButtonData(p.Name, fmt.Sprintf("product_%d", p.ID))
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func variantsKeyboard(product Product) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, v := range product.Variants {
		label := fmt.Sprintf("%s - %.2f₽", v.Color, v.Price)
		if !v.InStock {
			label += " (нет в наличии)"
		}
		button := tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("variant_%d_%d", product.ID, i))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
