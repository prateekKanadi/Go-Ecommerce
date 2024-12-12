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

func (repo *ProductRepository) getProduct(productID int) (*Product, error) {
	product := &Product{}
	whereClause := fmt.Sprintf("%s = ?", PRODUCT_ID)
	query := utils.BuildSelectQuery(TABLE_NAME, product, whereClause)

	row := repo.db.QueryRow(query, productID)
	err := row.Scan(
		&product.ProductID,
		&product.PricePerUnit,
		&product.ProductName,
		&product.ProductBrand,
		&product.Description,
		&product.StockQuantity,
		&product.Category,
		&product.SubCategory,
		&product.ImageURL)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return product, nil
}

func (repo *ProductRepository) getAllProducts() ([]Product, error) {
	query := utils.BuildSelectQuery(TABLE_NAME, &Product{}, "")
	results, err := repo.db.Query(query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.PricePerUnit,
			&product.ProductName,
			&product.ProductBrand,
			&product.Description,
			&product.StockQuantity,
			&product.Category,
			&product.SubCategory,
			&product.ImageURL)

		products = append(products, product)
	}
	return products, nil
}

func (repo *ProductRepository) getAllSimilarProducts(product *Product) ([]Product, error) {
	whereClause := fmt.Sprintf("category = ? AND %s != ?", PRODUCT_ID)
	query := utils.BuildSelectQuery(TABLE_NAME, &Product{}, whereClause)
	results, err := repo.db.Query(query, product.Category, product.ProductID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(&product.ProductID,
			&product.PricePerUnit,
			&product.ProductName,
			&product.ProductBrand,
			&product.Description,
			&product.StockQuantity,
			&product.Category,
			&product.SubCategory,
			&product.ImageURL)

		products = append(products, product)
	}
	return products, nil
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
