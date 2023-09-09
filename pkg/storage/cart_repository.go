package storage

import (
	"errors"
	"fmt"
	"trafilea-tech-challenge/pkg/models"
)

type CartRepository interface {
	CreateCart(userID string, cart models.Cart) models.Cart
	AddProduct(cartID string, product models.Product) (models.Cart, error)
	UpdateProductQuantity(cartID, product string, quantity int) (models.Cart, error)
	GetCartByID(cartID string) (models.Cart, error)
}

type cartRepo struct {
	repo map[string]models.Cart
}

func NewCartRepo(repo map[string]models.Cart) CartRepository {
	return &cartRepo{
		repo: repo,
	}
}

func (c *cartRepo) GetCartByID(cartID string) (models.Cart, error) {
	var cartToReturn models.Cart
	var found bool
	for _, cart := range c.repo {
		if cart.ID == cartID {
			cartToReturn = cart
			found = true
		}
	}

	if !found {
		return models.Cart{}, errors.New(fmt.Sprintf("cart with ID %v doesn't exist", cartID))
	}

	return cartToReturn, nil
}

func (c *cartRepo) UpdateProductQuantity(cartID, product string, quantity int) (models.Cart, error) {
	existingCart, err := c.GetCartByID(cartID)
	if err != nil {
		return models.Cart{}, errors.New(err.Error())
	}

	productInCart := findProductInCart(existingCart, product)
	if productInCart == nil {
		return models.Cart{}, errors.New(fmt.Sprintf("product %v does not exist in cart", product))
	}

	for i := 1; i < quantity; i++ {
		_, err := c.AddProduct(cartID, *productInCart)
		if err != nil {
			return models.Cart{}, err
		}
	}

	return c.repo[existingCart.UserID], nil
}

func (c *cartRepo) CreateCart(userID string, cartToCreate models.Cart) models.Cart {
	existingCart, ok := c.repo[userID]
	if ok {
		return existingCart
	}

	c.repo[userID] = cartToCreate
	return c.repo[userID]
}

func (c *cartRepo) AddProduct(cartID string, product models.Product) (models.Cart, error) {
	var userCart models.Cart
	var found bool
	for _, cart := range c.repo {
		if cart.ID == cartID {
			userCart = cart
			found = true
		}
	}

	if !found {
		return models.Cart{}, errors.New(fmt.Sprintf("cart with ID %v doesn't exist", cartID))
	}

	userCart.Products = append(userCart.Products, product)
	c.repo[userCart.UserID] = userCart
	return c.repo[userCart.UserID], nil
}

func findProductInCart(cart models.Cart, productToFind string) *models.Product {
	productsInCart := cart.Products
	for _, prod := range productsInCart {
		if prod.Name == productToFind {
			return &prod
		}
	}

	return nil
}
