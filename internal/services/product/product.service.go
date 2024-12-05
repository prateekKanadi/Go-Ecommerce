package product

import (
	"errors"
	"log"
	"net/http"
)

// ProductService handles business logic for product-related operations.
type ProductService struct {
	Repo *ProductRepository
}

// NewProductService creates a new ProductService.
func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}

func (s *ProductService) getAllProductsService() ([]Product, int, error) {
	productList, err := s.Repo.getAllProducts()
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return productList, http.StatusOK, nil
}

func (s *ProductService) GetProductService(productID int) (*Product, int, error) {
	product, err := s.Repo.getProduct(productID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if product == nil {
		return nil, http.StatusNotFound, errors.New("No Product Found")
	}

	return product, http.StatusOK, nil
}

func (s *ProductService) addProductService(newProduct Product) (int, error) {
	_, err := s.Repo.addProduct(newProduct)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *ProductService) updateProductService(updatedProduct Product) (int, error) {
	err := s.Repo.updateProduct(updatedProduct)

	if err != nil {
		log.Println(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *ProductService) removeProductService(productID int) (int, error) {
	err := s.Repo.removeProduct(productID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
