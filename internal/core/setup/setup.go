package setup

import (
	"log"

	"github.com/ecommerce/configuration"
	"github.com/ecommerce/database"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/sessions"
)

type InitializationResult struct {
	Config *configuration.Config // configuration type
	Store  *sessions.CookieStore // Or the exact type of your session store
}

// initialize all core components (except routes)
func InitializeAll(configPath string) (*InitializationResult, error) {
	result := &InitializationResult{}
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
	result.Config = config

	// Setup database
	err = database.SetupDatabase(config)
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return nil, err
	}

	// Setup session
	store, err := session.Init(config)
	if err != nil {
		log.Printf("Failed to initialize session: %v", err)
		return nil, err
	}
	result.Store = store

	return result, nil
}
