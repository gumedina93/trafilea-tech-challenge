package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"trafilea-tech-challenge/handlers"
	"trafilea-tech-challenge/pkg/cart"
	"trafilea-tech-challenge/pkg/models"
	"trafilea-tech-challenge/pkg/storage"
)

func main() {
	// Map used as in memory storage. For this example, we assume that one user can have only one cart
	var localStorage = make(map[string]models.Cart)

	cartRepo := storage.NewCartRepo(localStorage)
	cartService := cart.NewCart(cartRepo)

	router := gin.Default()
	router.POST("/carts", handlers.CreateCartHandler(cartService))
	router.POST("/carts/:cart_id/products", handlers.AddProductToCartHandler(cartService))
	router.PUT("/carts/:cart_id/products/:product", handlers.UpdateProductQuantityInCart(cartService))
	router.POST("/carts/:cart_id/orders", handlers.CreateOrderForCart(cartService))

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
