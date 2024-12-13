package product

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ecommerce/utils"
)

const (
	PRODUCT_ID = "productId"
	TABLE_NAME = "products"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}
func (repo *ProductRepository) getProduct(productID int, currency string) (*Product, error) {
	product := &Product{}

	// Ensure the currency is one of the valid options
	if currency != "USD" && currency != "EUR" && currency != "GBP" {
		return nil, fmt.Errorf("invalid currency: %s", currency)
	}

	// Build the SQL query dynamically based on the currency
	var priceColumn string
	switch currency {
	case "USD":
		priceColumn = "pp.price_usd"
	case "EUR":
		priceColumn = "pp.price_eur"
	case "GBP":
		priceColumn = "pp.price_gbp"
	}

	// Join products table with productPrices table to fetch the price based on the selected currency
	query := fmt.Sprintf(`
		SELECT p.productId, p.productName, p.productBrand, p.description, p.stockQuantity, %s
		FROM products p
		JOIN productPrices pp ON p.productId = pp.productId
		WHERE p.productId = ?
	`, priceColumn)

	// Fetch the product data from the database
	row := repo.db.QueryRow(query, productID)
	err := row.Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductBrand,
		&product.Description,
		&product.StockQuantity,
		&product.PricePerUnit, // Only the price for the selected currency is scanned here
	)
	if err == sql.ErrNoRows {
		return nil, nil // No product found
	} else if err != nil {
		log.Println("Error fetching product:", err)
		return nil, err
	}

	// Depending on the currency, set the correct price
	switch currency {
	case "USD":
		product.PricePerUnitUSD = product.PricePerUnit
	case "EUR":
		product.PricePerUnitEUR = product.PricePerUnit
	case "GBP":
		product.PricePerUnitGBP = product.PricePerUnit
	}

	log.Println("Product retrieved successfully:", product)
	return product, nil
}

func (repo *ProductRepository) removeProduct(productID int) error {
	whereClause := fmt.Sprintf("%s = ?", PRODUCT_ID)
	query := utils.BuildDeleteQuery(TABLE_NAME, whereClause)

	_, err := repo.db.Exec(query, productID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (repo *ProductRepository) getAllProducts(currency string) ([]Product, error) {

	// Ensure the currency is one of the valid options
	if currency != "USD" && currency != "EUR" && currency != "GBP" {
		return nil, fmt.Errorf("invalid currency: %s", currency)
	}

	// Build the SQL query dynamically based on the currency
	var priceColumn string
	switch currency {
	case "USD":
		priceColumn = "pp.price_usd"
	case "EUR":
		priceColumn = "pp.price_eur"
	case "GBP":
		priceColumn = "pp.price_gbp"
	}

	query := fmt.Sprintf(`
        SELECT p.productId, p.productName, p.productBrand, p.description, p.stockQuantity, %s
        FROM products p
        JOIN productPrices pp ON p.productId = pp.productId
    `, priceColumn)

	results, err := repo.db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)

	for results.Next() {
		var product Product
		// Scan the product fields and the selected price field based on currency
		err := results.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductBrand,
			&product.Description,
			&product.StockQuantity,
			&product.PricePerUnit, // The PricePerUnit field for the selected currency
		)
		if err != nil {
			log.Println("Error scanning row: ", err.Error())
			return nil, err
		}

		// Depending on the currency, set the correct price
		switch currency {
		case "USD":
			product.PricePerUnitUSD = product.PricePerUnit
		case "EUR":
			product.PricePerUnitEUR = product.PricePerUnit
		case "GBP":
			product.PricePerUnitGBP = product.PricePerUnit
		}

		// Append the product to the list
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepository) updateProduct(product Product) error {
	whereClause := fmt.Sprintf("%s = %d", PRODUCT_ID, product.ProductID)
	query, args := utils.BuildUpdateQuery(TABLE_NAME, product, whereClause)

	// Log the query and arguments to inspect them
	log.Println("Query:", query)
	log.Println("Args:", fmt.Sprintln(args...))

	_, err := repo.db.Exec(query, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (repo *ProductRepository) addProduct(product Product) (int, error) {
	query, args := utils.BuildInsertQuery(TABLE_NAME, product)
	result, err := repo.db.Exec(query, args...)

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
