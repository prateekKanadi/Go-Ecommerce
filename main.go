package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	// Create a server
	srv := &http.Server{
		Addr:    ":5000",
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exiting")
}
