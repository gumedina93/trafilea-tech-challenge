package cart

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"trafilea-tech-challenge/pkg/models"
	"trafilea-tech-challenge/pkg/storage"
)

func TestCreateCart_Success(t *testing.T) {
	// Given
	userID := "12345"
	testCart := models.Cart{
		ID:       mock.Anything,
		UserID:   userID,
		Products: []models.Product{},
	}

	repo := &storage.CartRepositoryMock{}
	repo.On("CreateCart", userID, mock.Anything).Return(testCart)
	cartService := NewCart(repo)

	// When
	userCart := cartService.CreateCart(userID)

	// Then
	require.Equal(t, userCart, testCart)
}

func TestAddProductToCart_Success_Extra_Coffee(t *testing.T) {
	// Given
	cartID := "test_cart_id"
	userID := "12345"

	coffeeProd := models.Product{
		Name:     "coffee2",
		Category: models.CoffeeCategory,
		Price:    20,
	}

	testCart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Products: []models.Product{
			{
				Name:     "coffee1",
				Category: models.CoffeeCategory,
				Price:    10,
			},
			{
				Name:     "coffee2",
				Category: models.CoffeeCategory,
				Price:    20,
			},
		},
	}

	repo := &storage.CartRepositoryMock{}
	repo.On("AddProduct", cartID, coffeeProd).Return(testCart, nil)

	extraCoffee := models.Product{
		Name:     "extraCoffee",
		Category: models.CoffeeCategory,
		Price:    0,
	}

	updatedTestCart := testCart
	updatedTestCart.Products = append(updatedTestCart.Products, extraCoffee)

	repo.On("AddProduct", cartID, extraCoffee).Return(updatedTestCart, nil)
	cartService := NewCart(repo)

	// When
	updatedCart, err := cartService.AddProductToCart(cartID, coffeeProd)

	// Then
	require.NoError(t, err)
	require.Equal(t, cartID, updatedCart.ID)
	require.Equal(t, userID, updatedCart.UserID)
	require.Equal(t, 3, len(updatedCart.Products))
}

func TestCreateOrderForCart_Error_Getting_Cart(t *testing.T) {
	// Given
	cartID := "test_cart_id"
	repo := &storage.CartRepositoryMock{}
	repo.On("GetCartByID", cartID).Return(models.Cart{}, errors.New("cart does not exist"))
	cartService := NewCart(repo)

	// When
	order, err := cartService.CreateOrderForCart(cartID)

	// Then
	require.Error(t, err)
	require.Equal(t, "cart does not exist", err.Error())
	require.Equal(t, models.Order{}, order)
}

func TestCreateOrderForCart_Success(t *testing.T) {
	// Given
	cartID := "test_cart_id"
	userID := "12345"
	testCart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Products: []models.Product{
			{
				Name:     "coffee1",
				Category: models.CoffeeCategory,
				Price:    10,
			},
			{
				Name:     "eq1",
				Category: models.EquipmentCategory,
				Price:    20,
			},
		},
	}
	repo := &storage.CartRepositoryMock{}
	repo.On("GetCartByID", cartID).Return(testCart, nil)
	cartService := NewCart(repo)

	// When
	order, err := cartService.CreateOrderForCart(testCart.ID)

	// Then
	require.NoError(t, err)
	require.Equal(t, cartID, order.CartID)
	require.Equal(t, 2, order.Totals.Products)
	require.Equal(t, fixedShippingPrice, order.Totals.Shipping)
	require.Equal(t, 0, order.Totals.Discounts)
}

func TestCreateOrderForCart_Success_With_Discounts(t *testing.T) {
	// Given
	cartID := "test_cart_id"
	userID := "12345"
	testCart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Products: []models.Product{
			{
				Name:     "acc1",
				Category: models.AccessoriesCategory,
				Price:    80,
			},
			{
				Name:     "eq1",
				Category: models.EquipmentCategory,
				Price:    20,
			},
			{
				Name:     "eq2",
				Category: models.EquipmentCategory,
				Price:    30,
			},
			{
				Name:     "eq3",
				Category: models.EquipmentCategory,
				Price:    20,
			},
			{
				Name:     "eq4",
				Category: models.EquipmentCategory,
				Price:    50,
			},
		},
	}
	repo := &storage.CartRepositoryMock{}
	repo.On("GetCartByID", cartID).Return(testCart, nil)
	cartService := NewCart(repo)

	// When
	order, err := cartService.CreateOrderForCart(testCart.ID)

	// Then
	require.NoError(t, err)
	require.Equal(t, cartID, order.CartID)
	require.Equal(t, 5, order.Totals.Products)
	require.Equal(t, 0, order.Totals.Shipping)
	require.Equal(t, 20, order.Totals.Discounts)
	require.Equal(t, 180, order.Totals.Price)
}

func TestUpdateProductQuantity_Success(t *testing.T) {
	// Given
	cartID := "test_cart_id"
	userID := "12345"
	testCart := models.Cart{
		ID:     cartID,
		UserID: userID,
		Products: []models.Product{
			{
				Name:     "acc1",
				Category: models.AccessoriesCategory,
				Price:    80,
			},
			{
				Name:     "coffee1",
				Category: models.CoffeeCategory,
				Price:    20,
			},
			{
				Name:     "coffee1",
				Category: models.CoffeeCategory,
				Price:    20,
			},
		},
	}
	repo := &storage.CartRepositoryMock{}
	repo.On("UpdateProductQuantity", cartID, "coffee1", 2).Return(testCart, nil)

	extraCoffee := models.Product{
		Name:     "extraCoffee",
		Category: models.CoffeeCategory,
		Price:    0,
	}

	updatedTestCart := testCart
	updatedTestCart.Products = append(updatedTestCart.Products, extraCoffee)
	repo.On("AddProduct", cartID, extraCoffee).Return(updatedTestCart, nil)
	cartService := NewCart(repo)

	// When
	userCart, err := cartService.UpdateProductQuantity(cartID, "coffee1", 2)

	// Then
	require.NoError(t, err)
	require.Equal(t, cartID, userCart.ID)
	require.Equal(t, userID, userCart.UserID)
	require.Equal(t, 4, len(userCart.Products))
}
