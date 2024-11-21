package routes

import (
	"database/sql"

	"github.com/ecommerce/internal/core/setup"
	"github.com/ecommerce/internal/services/authentication"
	"github.com/ecommerce/internal/services/index"
	"github.com/ecommerce/internal/services/product"
	"github.com/ecommerce/internal/services/user"
	"github.com/gorilla/mux"
)

type RepoInitializationResult struct {
	UserService    *user.UserService           // user service
	AuthService    *authentication.AuthService // auth service
	ProductService *product.ProductService     // prod service

}

func RegisterRoutes(r *mux.Router, setupRes *setup.InitializationResult) {
	//register repositories
	result := RegisterRepositories(setupRes.DbConn)

	index.SetupIndexRoutes(r)
	product.SetupProductRoutes(r, result.ProductService)
	user.SetupUserRoutes(r, result.UserService)
	authentication.SetupAuthRoutes(r, result.AuthService)
}

func RegisterRepositories(db *sql.DB) *RepoInitializationResult {
	result := &RepoInitializationResult{}

	//user repo/service initialize
	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	//auth repo/service initialize
	authService := authentication.NewAuthService(userService)

	//product repo/service initialize
	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo)

	result.UserService = userService
	result.AuthService = authService
	result.ProductService = productService

	return result
}
