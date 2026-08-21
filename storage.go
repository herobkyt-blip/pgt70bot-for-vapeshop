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
