package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/ecommerce/configuration"
	"github.com/ecommerce/internal/core/routes"
	"github.com/ecommerce/internal/core/setup"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	configFilePath = "config.yaml"
)

var (
	config *configuration.Config
	r      *mux.Router
)

func main() {
	fmt.Println("Hello! dev-anand")

	//setup configuration
	config, err := setup.InitializeAll(configFilePath)
	if err != nil {
		log.Println(config)
		log.Fatalf("Failed to initialize application: %v", err)
	}

	//Creating Mux Router
	r = mux.NewRouter()

	//Registering Middlewares
	setup.RegisterMiddleWares(r, config)

	//Registering routes
	routes.RegisterRoutes(r)

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
