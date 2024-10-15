package product

import (
	"database/sql"
	"log"

	"github.com/ecommerce/database"
)

func getProduct(productID int) (*Product, error) {
	row := database.DbConn.QueryRow(`SELECT 
	productId, 	
	pricePerUnit,	
	productName,
	productBrand
	FROM products 
	WHERE productId = ?`, productID)

	product := &Product{}
	err := row.Scan(
		&product.ProductID,
		&product.PricePerUnit,
		&product.ProductName,
		&product.ProductBrand)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return product, nil
}

func removeProduct(productID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM products where productId = ?`, productID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT 
	productId, 	 
	pricePerUnit, 	 
	productName,
	productBrand 
	FROM products`)
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
			&product.ProductBrand)

		products = append(products, product)
	}
	return products, nil
}

func updateProduct(product Product) error {
	_, err := database.DbConn.Exec(`UPDATE products SET 		 
		pricePerUnit=CAST(? AS DECIMAL(13,2)), 		 
		productName=?
		productBrand=?
		WHERE productId=?`,
		product.PricePerUnit,
		product.ProductName,
		product.ProductBrand,
		product.ProductID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`INSERT INTO products  
	(pricePerUnit,
	productName,
	productBrand) VALUES (?, ?, ?, ?, ?, ?)`,
		product.PricePerUnit,
		product.ProductName,
		product.ProductBrand)
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
