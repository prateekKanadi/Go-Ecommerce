package product

import (
	"errors"
	"log"
	"net/http"
)

func getAllProductsService() ([]Product, int, error) {
	productList, err := getAllProducts()
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return productList, http.StatusOK, nil
}

func getProductService(productID int) (*Product, int, error) {
	product, err := getProduct(productID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if product == nil {
		return nil, http.StatusNotFound, errors.New("No Product Found")
	}

	return product, http.StatusOK, nil
}

func addProductService(newProduct Product) (int, error) {
	_, err := addProduct(newProduct)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func updateProductService(updatedProduct Product) (int, error) {
	err := updateProduct(updatedProduct)

	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func removeProductService(productID int) (int, error) {
	err := removeProduct(productID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
