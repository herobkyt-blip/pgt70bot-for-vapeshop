package main

import "fmt"

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	InStock     bool
}

var products = []Product{
	{ID: 1, Name: "Жидкость TRAVA", Price: 550, Description: "Жидкость для электронных сигарет TRAVA", InStock: true},
	{ID: 2, Name: "Жидкость CHORME", Price: 540, Description: "Жидкость для электронных сигарет CHORME", InStock: true},
	{ID: 3, Name: "Жидкость PODONKI", Price: 450, Description: "Жидкость для электронных сигарет PODONKI", InStock: true},
}

func buildCatalogText() string {
	text := "Наш каталог:\n\n"
	for _, p := range products {
		text += fmt.Sprintf("%s - %.2f₽\n%s\n\n", p.Name, p.Price, p.Description)
	}
	return text
}
