package routes

import (
	"github.com/ecommerce/internal/core/services"
	"github.com/ecommerce/internal/core/setup"
	"github.com/ecommerce/internal/services/authentication"
	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/checkout"
	"github.com/ecommerce/internal/services/index"
	"github.com/ecommerce/internal/services/order"
	"github.com/ecommerce/internal/services/product"
	"github.com/ecommerce/internal/services/user"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, setupRes *setup.CoreSetupInitResult) {
	// Initialize services and repositories
	serviceRegistry := services.InitializeServices(setupRes.DbConn)

	// Register routes
	index.SetupIndexRoutes(r)
	product.SetupProductRoutes(r, serviceRegistry.ProductService)
	user.SetupUserRoutes(r, serviceRegistry.UserService)
	authentication.SetupAuthRoutes(r, serviceRegistry.AuthService)
	cart.SetupCartRoutes(r, serviceRegistry.CartService)
	checkout.SetupCheckoutRoutes(r, serviceRegistry.CheckoutService)
	order.SetupOrderRoutes(r, serviceRegistry.OrderService)
}
