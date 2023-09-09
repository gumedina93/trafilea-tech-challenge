package models

const (
	CoffeeCategory      = "coffee"
	EquipmentCategory   = "equipment"
	AccessoriesCategory = "accessories"
)

type Product struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Price    int    `json:"price"`
}

type Cart struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Products []Product `json:"products"`
}

type Order struct {
	CartID string `json:"cart_id"`
	Totals Total  `json:"totals"`
}

type Total struct {
	Products  int `json:"products"`
	Discounts int `json:"discounts"`
	Shipping  int `json:"shipping"`
	Order     int `json:"order"`
	Price     int `json:"price"`
}
