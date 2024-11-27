package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ecommerce/configuration"
)

// SetupDatabase
func SetupDatabase(config *configuration.Config) (*sql.DB, error) {
	var err error

	//  Build the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Database.User, config.Database.Password, config.Database.URL, config.Database.DbName)

	// Open a connection to the database
	dbConn, err := sql.Open("mysql", connStr) // root:root@tcp(127.0.0.1:3306)/ecommercedb
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}

	if err := dbConn.Ping(); err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	dbConn.SetMaxOpenConns(config.Database.MaxOpenConns)
	dbConn.SetMaxIdleConns(config.Database.MaxIdleConns)
	dbConn.SetConnMaxLifetime(time.Duration(config.Database.ConnMaxLifetime) * time.Second)

	return dbConn, nil
}
