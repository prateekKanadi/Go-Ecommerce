package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ecommerce/configuration"
)

var DbConn *sql.DB

// SetupDatabase
func SetupDatabase(config *configuration.Config) {
	var err error

	//  Build the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s", config.Database.User, config.Database.Password, config.Database.URL, config.Database.DbName)
	// log.Println(connStr)

	// Open a connection to the database
	DbConn, err = sql.Open("mysql", connStr) // root:root@tcp(127.0.0.1:3306)/ecommercedb
	if err != nil {
		log.Fatal(err)
	}

	DbConn.SetMaxOpenConns(3)
	DbConn.SetMaxIdleConns(3)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
