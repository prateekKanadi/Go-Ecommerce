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
	// Initialize product with an empty Prices slice
	product := &Product{
		Prices: make([]ProductPrice, 0),
	}

	// Ensure the currency is one of the valid options
	if currency != "USD" && currency != "EUR" && currency != "GBP" {
		return nil, fmt.Errorf("invalid currency: %s", currency)
	}

	// Query to fetch the product and its price for the specified currency
	query := `
		SELECT p.productId, p.productName, p.productBrand, p.description, p.stockQuantity, pp.price, pp.currencyCode, p.category, p.subCategory, p.imageURL
		FROM products p
		JOIN productPrices pp ON p.productId = pp.productId
		WHERE p.productId = ? AND pp.currencyCode = ?
	`

	// Fetch the product data for the selected product and currency
	row := repo.db.QueryRow(query, productID, currency)

	// Scan the result into individual fields
	var price float64
	var currencyCode string

	err := row.Scan(
		&product.ProductID,
		&product.ProductName,
		&product.ProductBrand,
		&product.Description,
		&product.StockQuantity,
		&price,        // Scan the price into a variable
		&currencyCode, // Scan the currencyCode into a variable
		&product.Category,
		&product.SubCategory, // This should be a string
		&product.ImageURL,
	)

	// Handle different error scenarios
	if err == sql.ErrNoRows {
		// No product found for the given productID and currency
		log.Printf("No product found for productID: %d and currency: %s", productID, currency)
		return nil, nil
	} else if err != nil {
		// Log the error if something else went wrong
		log.Printf("Error fetching product details: %v", err)
		return nil, fmt.Errorf("error fetching product details: %v", err)
	}

	// After scanning the price and currencyCode, append the price to the Prices slice
	product.Prices = append(product.Prices, ProductPrice{
		CurrencyCode: currencyCode,
		Amount:       price, // Assign the scanned price value
	})

	// Log the successful retrieval of the product
	log.Printf("Product retrieved successfully: %+v", product)

	// Return the product
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

	// Query to fetch all products and their price in the selected currency
	query := `
		SELECT p.productId, p.productName, p.productBrand, p.description, p.stockQuantity, pp.price, pp.currencyCode, 
		p.category, p.subCategory, p.imageURL
		FROM products p
		JOIN productPrices pp ON p.productId = pp.productId
		WHERE pp.currencyCode = ?
	`

	results, err := repo.db.Query(query, currency)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)

	// Iterate through each result and map it to a product
	for results.Next() {
		var product Product
		var price ProductPrice
		err := results.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ProductBrand,
			&product.Description,
			&product.StockQuantity,
			&price.Amount,       // The price for the selected currency
			&price.CurrencyCode, // The currency code
			&product.Category,
			&product.SubCategory,
			&product.ImageURL,
		)
		if err != nil {
			log.Println("Error scanning row: ", err.Error())
			return nil, err
		}

		// Add the price to the product's prices slice
		product.Prices = append(product.Prices, price)
		// Append the product to the list
		products = append(products, product)
	}

	return products, nil
}

func (repo *ProductRepository) getAllVariantProducts(currency string) ([]VariantProduct, error) {
	// Ensure the currency is one of the valid options
	if currency != "USD" && currency != "EUR" && currency != "GBP" {
		return nil, fmt.Errorf("invalid currency: %s", currency)
	}

	//  Query to fetch all variant products and their prices in the selected currency
	query := `
	SELECT 
		v.variantId, 
		v.variantName,
		v.description, 
		v.stockQuantity,
		v.imageURL,
		v.color,
		v.productId, 
		p.productName, 
		p.productBrand,
		pp.price, 
		pp.currencyCode, 
		p.category, 
		p.subCategory
	FROM 
		variantProducts v
	JOIN 
		products p ON v.productId = p.productId
	JOIN 
		productPrices pp ON p.productId = pp.productId
	WHERE 
		pp.currencyCode = ?
`

	results, err := repo.db.Query(query, currency)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()

	variantProducts := make([]VariantProduct, 0)

	// Iterate through each result and map it to a product
	for results.Next() {
		var product Product
		var price ProductPrice
		var variantProduct VariantProduct

		err := results.Scan(
			&variantProduct.VariantID,
			&variantProduct.VariantName,
			&variantProduct.Description,
			&variantProduct.StockQuantity,
			&variantProduct.ImageURL,
			&variantProduct.Color,
			&variantProduct.ProductID,
			&variantProduct.ProductName,
			&variantProduct.ProductBrand,
			&price.Amount,       // The price for the selected currency
			&price.CurrencyCode, // The currency code
			&variantProduct.Category,
			&variantProduct.SubCategory,
		)
		if err != nil {
			log.Println("Error scanning row: ", err.Error())
			return nil, err
		}

		// Add the price to the variant product's prices slice
		variantProduct.Prices = append(product.Prices, price)

		// Append the product to the list
		variantProducts = append(variantProducts, variantProduct)
	}

	return variantProducts, nil
}

func (repo *ProductRepository) getAllSimilarProducts(product *Product, currency string) ([]Product, error) {
	// Ensure the currency is one of the valid options
	if currency != "USD" && currency != "EUR" && currency != "GBP" {
		return nil, fmt.Errorf("invalid currency: %s", currency)
	}

	// Query to fetch all similar products and their price in the selected currency
	query := `
		SELECT p.productId, p.productName, p.productBrand, p.description, p.stockQuantity, pp.price, pp.currencyCode, 
		p.category, p.subCategory, p.imageURL
		FROM products p
		JOIN productPrices pp ON p.productId = pp.productId
		WHERE p.category = ? AND p.productId != ? AND pp.currencyCode = ?
	`

	results, err := repo.db.Query(query, product.Category, product.ProductID, currency)
	if err != nil {
		log.Println("Error fetching similar products: ", err)
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)

	// Iterate through each result and map it to a product
	for results.Next() {
		var similarProduct Product
		var price ProductPrice
		err := results.Scan(
			&similarProduct.ProductID,
			&similarProduct.ProductName,
			&similarProduct.ProductBrand,
			&similarProduct.Description,
			&similarProduct.StockQuantity,
			&price.Amount,       // The price for the selected currency
			&price.CurrencyCode, // The currency code
			&similarProduct.Category,
			&similarProduct.SubCategory,
			&similarProduct.ImageURL,
		)
		if err != nil {
			log.Println("Error scanning similar product row: ", err.Error())
			return nil, err
		}

		// Add the price to the product's Prices slice
		similarProduct.Prices = append(similarProduct.Prices, price)

		// Append the similar product to the products list
		products = append(products, similarProduct)
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

	// Add prices for this product in multiple currencies
	for _, price := range product.Prices {
		// Insert prices for different currencies into productPrices table
		priceQuery := `
			INSERT INTO productPrices (productId, price, currencyCode)
			VALUES (?, ?, ?)
		`
		_, err := repo.db.Exec(priceQuery, insertID, price.Amount, price.CurrencyCode)
		if err != nil {
			log.Println("Error inserting price:", err)
			return 0, err
		}
	}

	return int(insertID), nil
}
