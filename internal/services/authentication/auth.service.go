package authentication

import (
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/user"
)

// AuthService handles business logic for auth-related operations.
type AuthService struct {
	UserService *user.UserService
}

// NewAuthService creates a new AuthService.
func NewAuthService(userService *user.UserService) *AuthService {
	return &AuthService{
		UserService: userService,
	}
}

func (s *AuthService) registerUserService(newUser user.User) (int, error) {
	_, err := s.UserService.Repo.RegisterUser(newUser)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}

func (s *AuthService) loginUserService(existingUser user.User) (int, error) {
	res, err := s.UserService.Repo.LoginUser(existingUser)
	if err != nil {
		log.Print(err)
		return res, err
	}
	return res, nil
}
