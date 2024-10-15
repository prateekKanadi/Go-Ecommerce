package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ecommerce/database"
	"github.com/ecommerce/product"
)

const apiBasePath = "/api"

func main() {
	fmt.Println("Hello! dev-anand")

	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)

	fmt.Println("Server is running at http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
