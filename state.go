package main

type CartItem struct {
	ProductName string
	Color       string
	Price       float64
}

var carts = map[int64][]CartItem{}

func addToCart(userID int64, item CartItem) {
	carts[userID] = append(carts[userID], item)
}

func getCart(userID int64) []CartItem {
	return carts[userID]
}

type AdminStep string

const (
	StepNone                AdminStep = ""
	StepWaitingName         AdminStep = "waiting_name"
	StepWaitingPrice        AdminStep = "waiting_price"
	StepWaitingDesc         AdminStep = "waiting_description"
	StepWaitingCategory     AdminStep = "waiting_category"
	StepWaitingColor        AdminStep = "waiting_color"
	StepWaitingCategoryName AdminStep = "waiting_category_name"
)

type ProductDraft struct {
	Name        string
	Price       float64
	Description string
	Category    string
	Color       string
	Step        AdminStep
}

var drafts = map[int64]*ProductDraft{}
