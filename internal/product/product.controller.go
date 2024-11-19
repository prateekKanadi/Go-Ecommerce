package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	productsBasePath = "products"
	apiVersion       = "prod"
	apiBasePath      = "api"
	userRole         = "user"
)

// SetupRoutes :
func SetupProductRoutes(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", apiBasePath, productsBasePath), productsHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/{id}", apiBasePath, productsBasePath), productHandler)

	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("%s/%s", apiVersion, productsBasePath)
	r.HandleFunc(fmt.Sprintf("/%s/{id}", prodUrlPath), productProdHandler)
	r.HandleFunc(fmt.Sprintf("/%s", prodUrlPath), productsProdHandler)
}

func productsProdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("template/product_list.html")
		if err != nil {
			log.Println("Template parsing error:", err)
			http.Error(w, "Error loading product list page", http.StatusInternalServerError)
			return
		}

		productList, res, err := getAllProductsService()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		// Sample user role check
		// userRole := "user" // Replace with actual user role retrieval logic
		isAdmin := userRole == "admin"

		err = tmpl.Execute(w, map[string]interface{}{"Products": productList, "IsAdmin": isAdmin})
		if err != nil {
			log.Println("Template execution error:", err)
			http.Error(w, "Error rendering product list page", http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		// add a new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if newProduct.ProductID != 0 {
			err := errors.New("ProductId cannot be zero")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := addProductService(newProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productProdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	product, res, err := getProductService(productID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), res)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Parse the product details template
		tmpl, err := template.ParseFiles("template/product_details.html")
		if err != nil {
			log.Println("Template parsing error:", err)
			http.Error(w, "Error loading product details page", http.StatusInternalServerError)
			return
		}

		// Sample user role check
		// userRole := "user" // Replace with actual user role retrieval logic
		isAdmin := userRole == "admin"

		err = tmpl.Execute(w, map[string]interface{}{"Product": product, "IsAdmin": isAdmin})
		if err != nil {
			log.Println("Template execution error:", err)
			http.Error(w, "Error rendering product details page", http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPut:
		//update product in the list
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedProduct.ProductID != productID {
			err := errors.New("Payload Product Id Mismatch")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err = updateProductService(updatedProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		res, err := removeProductService(productID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, res, err := getAllProductsService()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		productsJson, err := json.Marshal(productList)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
		return
	case http.MethodPost:
		// add a new product to the list
		var newProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if newProduct.ProductID != 0 {
			err := errors.New("ProductId cannot be zero")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := addProductService(newProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	product, res, err := getProductService(productID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), res)
		return
	}

	switch r.Method {
	case http.MethodGet:
		//return single product
		productJson, err := json.Marshal(product)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
		return
	case http.MethodPut:
		//update product in the list
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedProduct.ProductID != productID {
			err := errors.New("Payload Product Id Mismatch")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err = updateProductService(updatedProduct)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		res, err := removeProductService(productID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
