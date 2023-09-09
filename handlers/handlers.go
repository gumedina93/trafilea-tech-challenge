package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trafilea-tech-challenge/pkg/cart"
	"trafilea-tech-challenge/pkg/models"
)

func CreateOrderForCart(cartService cart.Cart) gin.HandlerFunc {
	return func(c *gin.Context) {
		cartID := c.Param("cart_id")
		order, err := cartService.CreateOrderForCart(cartID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func UpdateProductQuantityInCart(cartService cart.Cart) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			NewQuantity int `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if request.NewQuantity < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product quantity must be greater than 0"})
			return
		}

		product := c.Param("product")
		cartID := c.Param("cart_id")
		updatedCart, err := cartService.UpdateProductQuantity(cartID, product, request.NewQuantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedCart)
	}
}

func CreateCartHandler(cartService cart.Cart) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			UserID string `json:"user_id"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userCart := cartService.CreateCart(request.UserID)
		c.JSON(http.StatusOK, userCart)
	}
}

func AddProductToCartHandler(cartService cart.Cart) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Name     string `json:"name"`
			Category string `json:"category"`
			Price    int    `json:"price"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !isValidCategory(request.Category) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
			return
		}

		if request.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no empty values allowed"})
			return
		}

		if request.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "price must be greater than 0"})
			return
		}

		cartID := c.Param("cart_id")
		product := models.Product{
			Name:     request.Name,
			Category: request.Category,
			Price:    request.Price,
		}

		res, err := cartService.AddProductToCart(cartID, product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func isValidCategory(category string) bool {
	return category == models.CoffeeCategory || category == models.EquipmentCategory || category == models.AccessoriesCategory
}
