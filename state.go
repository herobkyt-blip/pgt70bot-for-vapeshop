package main

var carts = map[int64][]Product{}

func addToCart(userID int64, product Product) {
	carts[userID] = append(carts[userID], product)
}

func getCart(userID int64) []Product {
	return carts[userID]
}
