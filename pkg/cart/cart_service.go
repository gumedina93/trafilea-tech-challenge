package cart

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
	"trafilea-tech-challenge/pkg/models"
	"trafilea-tech-challenge/pkg/storage"
)

const (
	fixedShippingPrice = 20
)

type Cart interface {
	CreateCart(userID string) models.Cart
	AddProductToCart(cartID string, product models.Product) (models.Cart, error)
	UpdateProductQuantity(cartID, product string, quantity int) (models.Cart, error)
	CreateOrderForCart(cartID string) (models.Order, error)
}

type cart struct {
	CartRepo storage.CartRepository
}

func NewCart(storage storage.CartRepository) Cart {
	return &cart{
		CartRepo: storage,
	}
}

func (c *cart) CreateOrderForCart(cartID string) (models.Order, error) {
	userCart, err := c.CartRepo.GetCartByID(cartID)
	if err != nil {
		return models.Order{}, err
	}

	orderID := generateOrderID()
	order := models.Order{
		CartID: cartID,
		Totals: models.Total{
			Order:    orderID,
			Shipping: fixedShippingPrice,
		},
	}

	productsQuantityByCategory := getProductsQuantityByCategory(userCart)
	if productsQuantityByCategory.Equipment > 3 {
		order.Totals.Shipping = 0
	}

	totalSpent, totalProducts, discount := calculateOrderDetails(userCart)
	order.Totals.Price = totalSpent
	order.Totals.Products = totalProducts
	order.Totals.Discounts = discount

	return order, nil
}

func (c *cart) UpdateProductQuantity(cartID, product string, quantity int) (models.Cart, error) {
	updatedCart, err := c.CartRepo.UpdateProductQuantity(cartID, product, quantity)
	if err != nil {
		return models.Cart{}, err
	}

	productsQuantityByCategory := getProductsQuantityByCategory(updatedCart)
	productCategory := getProductCategoryByName(updatedCart, product)
	hasFreeCoffee := true
	if productCategory == models.CoffeeCategory {
		hasFreeCoffee = hasAlreadyFreeCoffee(updatedCart)
	}
	if productsQuantityByCategory.Coffee >= 2 && !hasFreeCoffee {
		updatedCart, err = c.CartRepo.AddProduct(cartID, models.Product{
			Name:     "extraCoffee",
			Category: models.CoffeeCategory,
			Price:    0,
		})

		if err != nil {
			return models.Cart{}, err
		}
	}

	return updatedCart, nil
}

func (c *cart) AddProductToCart(cartID string, product models.Product) (models.Cart, error) {
	updatedCart, err := c.CartRepo.AddProduct(cartID, product)
	if err != nil {
		return models.Cart{}, err
	}

	productsQuantityByCategory := getProductsQuantityByCategory(updatedCart)
	hasFreeCoffee := hasAlreadyFreeCoffee(updatedCart)
	if productsQuantityByCategory.Coffee >= 2 && !hasFreeCoffee {
		updatedCart, err = c.CartRepo.AddProduct(cartID, models.Product{
			Name:     "extraCoffee",
			Category: models.CoffeeCategory,
			Price:    0,
		})

		if err != nil {
			return models.Cart{}, err
		}
	}

	return updatedCart, nil
}

func (c *cart) CreateCart(userID string) models.Cart {
	newCart := models.Cart{
		ID:       uuid.New().String(),
		UserID:   userID,
		Products: []models.Product{},
	}

	return c.CartRepo.CreateCart(userID, newCart)
}

func calculateOrderDetails(cart models.Cart) (int, int, int) {
	categoryTotals := make(map[string]int)
	totalSpent := 0

	for _, product := range cart.Products {
		category := product.Category
		price := product.Price

		if _, exists := categoryTotals[category]; !exists {
			categoryTotals[category] = 0
		}

		categoryTotals[category] += price
		totalSpent += price
	}

	var discount int

	if categoryTotals[models.AccessoriesCategory] > 70 {
		discount = int(float64(totalSpent) * 0.10)
		totalSpent -= discount
	}

	return totalSpent, len(cart.Products), discount
}

func getProductCategoryByName(cart models.Cart, productName string) string {
	for _, product := range cart.Products {
		if product.Name == productName {
			return product.Category
		}
	}

	return ""
}

func hasAlreadyFreeCoffee(cart models.Cart) bool {
	for _, product := range cart.Products {
		if product.Category == models.CoffeeCategory && product.Price == 0 {
			return true
		}
	}

	return false
}

func getProductsQuantityByCategory(cart models.Cart) productsByCategory {
	prodsByCategory := productsByCategory{}
	for _, prod := range cart.Products {
		switch prod.Category {
		case models.CoffeeCategory:
			prodsByCategory.Coffee++
		case models.EquipmentCategory:
			prodsByCategory.Equipment++
		case models.AccessoriesCategory:
			prodsByCategory.Accessories++
		}
	}

	return prodsByCategory
}

func generateOrderID() int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	digits := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	r.Shuffle(len(digits), func(i, j int) {
		digits[i], digits[j] = digits[j], digits[i]
	})

	result := 0
	for i := 0; i < 8; i++ {
		result = result*10 + digits[i]
	}

	return result
}

type productsByCategory struct {
	Coffee      int
	Equipment   int
	Accessories int
}
