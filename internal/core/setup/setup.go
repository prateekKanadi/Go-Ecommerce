package setup

import (
	"database/sql"
	"log"

	"github.com/ecommerce/configuration"
	"github.com/ecommerce/database"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/sessions"
)

type CoreSetupInitResult struct {
	Config *configuration.Config // configuration type
	Store  *sessions.CookieStore // Or the exact type of your session store
	DbConn *sql.DB               // Database type
}

// initialize all core components (except routes)
func InitializeAll(configPath string) (*CoreSetupInitResult, error) {
	result := &CoreSetupInitResult{}

	// Setup yaml configuration
	config, err := configuration.Init(configPath)
	if err != nil {
		log.Printf("Failed to initialize configuration from %s: %v", configPath, err)
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		log.Printf("Failed to Validate initialized configuration from %s: %v", configPath, err)
		return nil, err
	}
	result.Config = config

	// Setup database
	dbConn, err := database.SetupDatabase(config)
	if err != nil {
		log.Printf("Failed to initialize database from %s: %v", configPath, err)
		return nil, err
	}
	result.DbConn = dbConn

	// Setup session
	store, err := session.Init(config)
	if err != nil {
		log.Printf("Failed to initialize session: %v from %s", configPath, err)
		return nil, err
	}
	result.Store = store

	return result, nil
}
