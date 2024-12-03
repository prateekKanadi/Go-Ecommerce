package cart

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/product"
)

// CartService handles business logic for product-related operations.
type CartService struct {
	Repo           *CartRepository
	ProductService *product.ProductService
}

// NewCartService creates a new CartService.
func NewCartService(repo *CartRepository, productService *product.ProductService) *CartService {
	return &CartService{
		Repo:           repo,
		ProductService: productService,
	}
}

// AddOrUpdateCartItem adds a product to the cart or updates the quantity if it already exists in the cart.
func (s *CartService) AddOrUpdateCartItemService(cartID, productID, quantity int) (int, error) {
	// Ensure the quantity is greater than zero
	if quantity <= 0 {
		return http.StatusBadRequest, fmt.Errorf("invalid quantity: must be greater than zero")
	}

	// Call the repository to perform the upsert
	err := s.Repo.addOrUpdateCartItem(cartID, productID, quantity)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to add or update cart item: %v", err)
	}

	return http.StatusOK, nil
}

/*Remove cart item method */
func (s *CartService) removeCartItem(cartID, productID int) (int, error) {

	err := s.Repo.removeCartItem(cartID, productID)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("failed to remove cart item: %v", err)
	}

	return http.StatusOK, nil
}

func (s *CartService) getAllCartItemsService(cartID int) (*Cart, int, error) {
	cartList, err := s.Repo.getAllCartItems(cartID)
	if err != nil {
		log.Printf("Error fetching cart items: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return cartList, http.StatusOK, nil
}
