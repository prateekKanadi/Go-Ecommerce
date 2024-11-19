package routes

import (
	"github.com/ecommerce/internal/services/authentication"
	"github.com/ecommerce/internal/services/index"
	"github.com/ecommerce/internal/services/product"
	"github.com/ecommerce/internal/services/user"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	index.SetupIndexRoutes(r)
	product.SetupProductRoutes(r)
	user.SetupUserRoutes(r)
	authentication.SetupAuthRoutes(r)
}
