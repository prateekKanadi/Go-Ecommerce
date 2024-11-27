package authentication

import (
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/user"
)

// AuthService handles business logic for auth-related operations.
type AuthService struct {
	UserService *user.UserService
	CartService *cart.CartService
}

// NewAuthService creates a new AuthService.
func NewAuthService(userService *user.UserService, cartService *cart.CartService) *AuthService {
	return &AuthService{
		UserService: userService,
		CartService: cartService,
	}
}

func (s *AuthService) registerUserService(newUser user.User) (int, int, error) {
	insertID, err := s.UserService.Repo.RegisterUser(newUser)
	if err != nil {
		log.Print(err)
		return 0, http.StatusBadRequest, nil
	}
	return insertID, http.StatusOK, err
}

func (s *AuthService) loginUserService(existingUser user.User) (int, error) {
	res, err := s.UserService.Repo.LoginUser(existingUser)
	if err != nil {
		log.Print(err)
		return res, err
	}
	return res, nil
}
