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

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

const (
	productsBasePath = "products"
	prodBasePath     = "prod"
	apiBasePath      = "api"
)

// SetupRoutes :
func SetupProductRoutes(r *mux.Router, s *ProductService) {
	apiUrlPath := fmt.Sprintf("/%s/%s", apiBasePath, productsBasePath)
	productRouter := r.PathPrefix(apiUrlPath).Subrouter()

	productRouter.HandleFunc("", productsHandler(s))
	productRouter.HandleFunc("/{id}", productHandler(s))

	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, productsBasePath)
	prodUsersRouter := r.PathPrefix(prodUrlPath).Subrouter()

	prodUsersRouter.HandleFunc("", productsProdHandler(s))
	prodUsersRouter.HandleFunc("/{id}", productProdHandler(s))
}

func productsProdHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//extracting isAnon flag from session
		isAnon := sess.Values["isAnon"].(bool)
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}
		if isAnon {
			log.Println("my products Anon userID : ", user.UserID)
		}

		switch r.Method {
		case http.MethodGet:
			tmpl, err := template.ParseFiles("template/product_list.html")
			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading product list page", http.StatusInternalServerError)
				return
			}
			currency := sess.Values["currency"].(string)
			productList, res, err := s.getAllProductsService(currency)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}

			err = tmpl.Execute(w, map[string]interface{}{"Products": productList, "IsAdmin": user.IsAdmin, "isAnon": isAnon})
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

			res, err := s.addProductService(newProduct)
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
}

func productProdHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//extracting isAnon flag from session
		isAnon := sess.Values["isAnon"].(bool)
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}
		if isAnon {
			log.Println("my product details Anon userID : ", user.UserID)
		}

		vars := mux.Vars(r)
		productID, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		currency := sess.Values["currency"].(string)
		product, res, err := s.GetProductService(productID, currency)
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

			err = tmpl.Execute(w, map[string]interface{}{"Product": product, "IsAdmin": user.IsAdmin, "isAnon": isAnon})
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

			res, err = s.updateProductService(updatedProduct)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			res, err := s.removeProductService(productID)
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
}

func productsHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch r.Method {
		case http.MethodGet:
			currency := sess.Values["currency"].(string)
			productList, res, err := s.getAllProductsService(currency)
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

			// adding product
			res, err := s.addProductService(newProduct)
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
}

func productHandler(s *ProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		productID, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		currency := sess.Values["currency"].(string)
		product, res, err := s.GetProductService(productID, currency)
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

			// update product cred
			res, err = s.updateProductService(updatedProduct)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			res, err := s.removeProductService(productID)
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
}
