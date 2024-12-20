package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce/internal/core/middleware"
	"github.com/ecommerce/internal/core/routes"
	"github.com/ecommerce/internal/core/setup"
	"github.com/ecommerce/internal/services/index"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	configFilePath = "config.yaml"
)

func main() {
	fmt.Println("Hello!")

	//setup configuration
	setupRes, err := setup.InitializeAll(configFilePath)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer setupRes.DbConn.Close() // Always close the connection when the application exits

	//Creating Mux Router
	r := mux.NewRouter()

	// Static file server for images
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./resources/images/"))))

	//Registering Middlewares
	middleware.RegisterMiddleWares(r, setupRes)

	//Registering routes
	routes.RegisterRoutes(r, setupRes)

	fmt.Println("Server is running at http://localhost:5000")

	// Automatically open the landing page in the default browser
	go index.ServeIndexPage()

	log.Fatal(http.ListenAndServe(":5000", r))
}
