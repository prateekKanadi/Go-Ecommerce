package setup

import (
	"log"

	"github.com/ecommerce/configuration"
	"github.com/ecommerce/database"
	"github.com/ecommerce/internal/core/middleware"
	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

// initialize all core components (except routes)
func InitializeAll(configPath string) (*configuration.Config, error) {
	// Setup yaml configuration
	config, err := configuration.Init(configPath)
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		log.Printf("Failed to Validate initialized configuration: %v", err)
		return nil, err
	}

	// Setup database
	err = database.SetupDatabase(config)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return nil, err
	}

	// Setup session
	// var store *sessions.CookieStore
	_, err = session.Init(config)
	if err != nil {
		log.Printf("Failed to initialize session: %v", err)
		return nil, err
	}

	return config, nil
}

// Setup middleware
func RegisterMiddleWares(r *mux.Router, config *configuration.Config) {
	r.Use(middleware.InjectConfigMiddleware(config)) // Add middleware for injecting config
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.SessionMiddleware(config))
}
