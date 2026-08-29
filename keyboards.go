package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func adminMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📁 Категории", "admin_categories_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📦 Товары", "admin_products_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🖼 Баннеры", "admin_banners_menu"),
		),
	)
}

func confirmVariantsKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить✅", "confirm_variants"),
		),
	)
}

func adminCategoriesMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕Добавить категорию", "admin_add_category"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить категорию", "admin_delete_category"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin_main_menu"),
		),
	)
}

func adminProductsMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕Добавить товар", "admin_add_product"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить товар", "admin_delete_product"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin_main_menu"),
		),
	)
}

func categoriesKeyboard(categoryList []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categoryList {
		button := tgbotapi.NewInlineKeyboardButtonData(category, "category_"+category)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	backButton := tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_main")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))
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
	backButton := tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_categories")
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))
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
	backButton := tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "back_products_"+product.Category)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(backButton))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func persistentKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 Каталог"),
			tgbotapi.NewKeyboardButton("🛒 Корзина"),
		),
	)
}

func adminCategoriesKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(category, "pickcat_"+category)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func deleteCategoryKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(category, "delcat_"+category)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func deleteProductKeyboad() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, p := range products {
		button := tgbotapi.NewInlineKeyboardButtonData(p.Name, fmt.Sprintf("delprod_%d", p.ID))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func checkoutKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Оформить заказ", "checkout"),
		),
	)
}

func cartKeyboard(cart []CartItem) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i, item := range cart {
		label := fmt.Sprintf("❌ %s (%s)", item.ProductName, item.Color)
		button := tgbotapi.NewInlineKeyboardButtonData(label, fmt.Sprintf("removecart_%d", i))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Оформить заказ", "checkout"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func bannersMenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👋 Приветствие", "setbanner_welcome"),
	))
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("📁 Каталог (список категорий)", "setbanner_catalog"),
	))
	for _, category := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(category, "setbanner_"+category)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "admin_main_menu"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
