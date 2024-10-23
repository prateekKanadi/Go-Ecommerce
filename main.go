package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce/authentication"
	"github.com/ecommerce/database"
	"github.com/ecommerce/middleware"
	"github.com/ecommerce/product"
	"github.com/ecommerce/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const apiBasePath = "/api"

func main() {
	fmt.Println("Hello! dev-anand")

	//Creating Mux Router
	r := mux.NewRouter()

	//register datatbase
	database.SetupDatabase()

	//Registering Middlewares
	registerMiddleWares(r)

	//Registering routes
	registerRoutes(r, apiBasePath)

	fmt.Println("Server is running at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func registerMiddleWares(r *mux.Router) {
	r.Use(middleware.CorsMiddleware)
}

func registerRoutes(r *mux.Router, apiBasePath string) {
	//Product
	product.SetupProductRoutes(r, apiBasePath)

	// User
	user.SetupUserRoutes(r, apiBasePath)

	// Auth
	authentication.SetupAuthRoutes(r, apiBasePath)
}
