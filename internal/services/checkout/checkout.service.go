package checkout

import (
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/user"
)

type CheckoutService struct{
	UserService *user.UserService
	CartService *cart.CartService
}

// Checkout Service
func NewCheckoutService(userService *user.UserService, cartService *cart.CartService) *CheckoutService {
	return &CheckoutService{
		UserService: userService,
		CartService: cartService,
	}
}

func (s *CheckoutService) getAddressDetailsOfUser(userId int) (*user.Address,int,error){
	addressDetails, err := s.UserService.Repo.GetAddressByUserId(userId)
	if err != nil {
		log.Printf("Error fetching user details: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return &addressDetails, http.StatusOK, nil
}

func (s *CheckoutService) getCartDetailsOfUser(cartId int) (*cart.Cart,error){
	cartData,err := s.CartService.Repo.GetAllCartItems(cartId)
	if err != nil {
		log.Printf("Error fetching cart details of user: %v", err)
		return nil, err
	}
	return cartData,nil
}