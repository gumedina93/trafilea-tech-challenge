package storage

import (
	"github.com/stretchr/testify/require"
	"testing"
	"trafilea-tech-challenge/pkg/models"
)

func TestCartRepo_GetCartByID(t *testing.T) {
	// Given
	repo := NewCartRepo(map[string]models.Cart{
		"12345": {
			ID:     "testCartID",
			UserID: "12345",
			Products: []models.Product{
				{Name: "product1", Category: models.CoffeeCategory, Price: 10},
			},
		},
	})

	// When
	cart, err := repo.GetCartByID("testCartID")

	// Then
	require.NoError(t, err)
	require.Equal(t, "12345", cart.UserID)
}

func TestCartRepo_UpdateProductQuantity(t *testing.T) {
	// Given
	repo := NewCartRepo(map[string]models.Cart{
		"12345": {
			ID:     "testCartID",
			UserID: "testUserID",
			Products: []models.Product{
				{Name: "product1", Category: models.CoffeeCategory, Price: 10},
			},
		},
	})

	// When
	updatedCart, err := repo.UpdateProductQuantity("testCartID", "product1", 3)

	// Then
	require.NoError(t, err)
	require.Equal(t, 3, len(updatedCart.Products))
}

func TestCartRepo_CreateCart_And_Add_Product(t *testing.T) {
	// Given
	repo := NewCartRepo(make(map[string]models.Cart))
	newCart := models.Cart{
		UserID: "testUserID",
	}

	createdCart := repo.CreateCart("testUserID", newCart)
	require.Equal(t, "testUserID", createdCart.UserID)

	// When
	res, err := repo.AddProduct(createdCart.ID, models.Product{
		Name:     "coffeeTest",
		Category: models.CoffeeCategory,
		Price:    15,
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, 1, len(res.Products))
}
