package main

type Variant struct {
	Color   string
	Price   float64
	InStock bool
}

type Product struct {
	ID             int
	Name           string
	Description    string
	Category       string
	PhotoFileID    string
	Variants       []Variant
	DiscountAmount float64
}

var products = []Product{}
var categories = []string{}
var banners = map[string]string{}

func findProduct(id int) (Product, bool) {
	for _, p := range products {
		if p.ID == id {
			return p, true

		}
	}
	return Product{}, false
}
