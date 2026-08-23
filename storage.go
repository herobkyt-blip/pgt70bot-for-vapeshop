package main

import (
	"encoding/json"
	"os"
)

const productFile = "products.json"

func saveProducts() {
	data, err := json.MarshalIndent(products, "", " ")
	if err != nil {
		return
	}
	os.WriteFile(productFile, data, 0644)
}

func loadProducts() {
	data, err := os.ReadFile(productFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &products)
}

const categoriesFile = "categories.json"

func saveCategories() {
	data, err := json.MarshalIndent(categories, "", " ")
	if err != nil {
		return
	}
	os.WriteFile(categoriesFile, data, 0644)
}

func loadCategories() {
	data, err := os.ReadFile(categoriesFile)
	if err != nil {
		return
	}
	json.Unmarshal(data, &categories)
}
