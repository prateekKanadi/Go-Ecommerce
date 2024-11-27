package services

import (
	"database/sql"

	"github.com/ecommerce/internal/services/authentication"
	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/product"
	"github.com/ecommerce/internal/services/user"
)

type ServiceRegistry struct {
	UserService    *user.UserService
	AuthService    *authentication.AuthService
	ProductService *product.ProductService
	CartService    *cart.CartService
}

func InitializeServices(db *sql.DB) *ServiceRegistry {
	// Initialize user repository and service
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	// Initialize cart repository and service
	cartRepo := cart.NewCartRepository(db)
	cartService := cart.NewCartService(cartRepo)

	// Initialize authentication service
	authService := authentication.NewAuthService(userService, cartService)

	// Initialize product repository and service
	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo)

	// Return the ServiceRegistry with all services initialized
	return &ServiceRegistry{
		UserService:    userService,
		AuthService:    authService,
		ProductService: productService,
		CartService:    cartService,
	}
}
