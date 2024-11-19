package routes

import (
	"github.com/ecommerce/internal/authentication"
	"github.com/ecommerce/internal/index"
	"github.com/ecommerce/internal/product"
	"github.com/ecommerce/internal/user"
	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	index.SetupIndexRoutes(r)
	product.SetupProductRoutes(r)
	user.SetupUserRoutes(r)
	authentication.SetupAuthRoutes(r)
}
