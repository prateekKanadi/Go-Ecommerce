package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/ecommerce/authentication"
	"github.com/ecommerce/configuration"
	"github.com/ecommerce/database"
	"github.com/ecommerce/index"
	"github.com/ecommerce/middleware"
	"github.com/ecommerce/product"
	"github.com/ecommerce/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	apiBasePath    = "/api"
	configFilePath = "config.yaml"
)

var (
	config *configuration.Config
	r      *mux.Router
)

func main() {
	fmt.Println("Hello! dev-anand")

	//setup configuration
	config = configuration.Init(configFilePath)

	//Creating Mux Router
	r = mux.NewRouter()

	//register datatbase
	database.SetupDatabase(config)

	//Registering Middlewares
	registerMiddleWares(r)

	//Registering routes
	registerRoutes(r, apiBasePath)

	fmt.Println("Server is running at http://localhost:5000")

	// Automatically open the landing page in the default browser
	go serveIndexPage()

	log.Fatal(http.ListenAndServe(":5000", r))
}

func serveIndexPage() {
	time.Sleep(1 * time.Second) // Wait a second for the server to start
	err := exec.Command("cmd", "/C", "start", "http://localhost:5000").Run()
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

func registerMiddleWares(r *mux.Router) {
	r.Use(middleware.CorsMiddleware)
	r.Use(middleware.AuthMiddleware)
}

func registerRoutes(r *mux.Router, apiBasePath string) {
	// Serve static files from the "static" directory
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//Landing Page
	index.SetupIndexRoutes(r)

	//Product
	product.SetupProductRoutes(r, apiBasePath)

	// User
	user.SetupUserRoutes(r, apiBasePath)

	// Auth
	authentication.SetupAuthRoutes(r, apiBasePath)
}
