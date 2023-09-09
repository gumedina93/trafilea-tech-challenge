package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"trafilea-tech-challenge/pkg/cart"
	"trafilea-tech-challenge/pkg/models"
)

func TestCreateCart_Success(t *testing.T) {
	// Given
	cartService := &cart.CartMock{}
	cartService.On("CreateCart", "123").Return(models.Cart{}, nil)

	r := gin.Default()
	r.POST("/carts", CreateCartHandler(cartService))
	reqBody := []byte(`{"user_id": "123"}`)
	req, err := http.NewRequest("POST", "/carts", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, err)
}

func TestCreateOrderForCart_Success(t *testing.T) {
	// Given
	cartService := &cart.CartMock{}
	cartService.On("CreateOrderForCart", "123").Return(models.Order{
		CartID: "123",
		Totals: models.Total{},
	}, nil)

	r := gin.Default()
	r.POST("/carts/:cart_id/orders", CreateOrderForCart(cartService))
	req, err := http.NewRequest("POST", "/carts/123/orders", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	var response models.Order
	err = json.Unmarshal(w.Body.Bytes(), &response)

	// Then
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, err)
	require.Equal(t, "123", response.CartID)
}

func TestUpdateProductQuantityInCart_Success(t *testing.T) {
	// Given
	cartService := &cart.CartMock{}
	cartService.On("UpdateProductQuantity", "1", "coffeeTest", 2).Return(models.Cart{
		ID:     "1",
		UserID: "19",
		Products: []models.Product{
			{
				Name:     "coffeeTest",
				Category: models.CoffeeCategory,
				Price:    10,
			},
			{
				Name:     "coffeeTest",
				Category: models.CoffeeCategory,
				Price:    10,
			},
		},
	}, nil)

	r := gin.Default()
	r.PUT("/carts/:cart_id/products/:product", UpdateProductQuantityInCart(cartService))
	reqBody := []byte(`{"quantity": 2}`)
	req, err := http.NewRequest("PUT", "/carts/1/products/coffeeTest", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAddProductToCart_Success(t *testing.T) {
	// Given
	cartService := &cart.CartMock{}

	product := models.Product{
		Name:     "coffeeA",
		Category: models.CoffeeCategory,
		Price:    15,
	}

	cartService.On("AddProductToCart", "1", product).Return(models.Cart{
		ID:     "1",
		UserID: "19",
		Products: []models.Product{
			product,
		},
	}, nil)

	r := gin.Default()
	r.POST("/carts/:cart_id/products", AddProductToCartHandler(cartService))
	reqBody := []byte(`{"name": "coffeeA", "category": "coffee", "price": 15}`)
	req, err := http.NewRequest("POST", "/carts/1/products", bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	w := httptest.NewRecorder()

	// When
	r.ServeHTTP(w, req)

	// Then
	var userCart models.Cart
	err = json.Unmarshal(w.Body.Bytes(), &userCart)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, 1, len(userCart.Products))
}
