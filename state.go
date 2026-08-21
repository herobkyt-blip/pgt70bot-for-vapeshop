package main

var carts = map[int64][]Product{}

func addToCart(userID int64, product Product) {
	carts[userID] = append(carts[userID], product)
}

func getCart(userID int64) []Product {
	return carts[userID]
}

type AdminStep string

const (
	StepNone         AdminStep = ""
	StepWaitingName  AdminStep = "waiting_name"
	StepWaitingPrice AdminStep = "waiting_price"
	StepWaitingDesc  AdminStep = "waiting_description"
)

type ProductDraft struct {
	Name        string
	Price       float64
	Description string
	Step        AdminStep
}

var drafts = map[int64]*ProductDraft{}
