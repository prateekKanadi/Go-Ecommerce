package cart

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

const (
	cartBasePath = "cart"
	prodBasePath = "prod"
	apiBasePath  = "api"
)

// SetupRoutes :
func SetupCartRoutes(r *mux.Router, s *CartService) {
	// apiUrlPath := fmt.Sprintf("/%s/%s", apiBasePath, cartBasePath)
	// cartRouter := r.PathPrefix(apiUrlPath).Subrouter()

	// cartRouter.HandleFunc("/{id}", AddToCartHandler(s))

	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, cartBasePath)
	prodCartRouter := r.PathPrefix(prodUrlPath).Subrouter()

	prodCartRouter.HandleFunc("/{id}", AddToCartProdHandler(s))
	// prodCartRouter.HandleFunc("/addCart/{id}", AddCartProdHandler(s))
}

// func AddCartProdHandler(s *CartService) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		UserID, err := strconv.Atoi(mux.Vars(r)["id"])
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}

// 		// Now, create a cart for the user
// 		cartID, err := s.Repo.CreateCartForUser(UserID)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		log.Println("inside add-cart-prod-handler cartID: ", cartID)

// 		// Success response
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Cart created successfully"))
// 	}
// }

func AddToCartProdHandler(s *CartService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cart := sess.Values["cart"].(*session.Cart)

		// Example values from request (you would extract these from form data, JSON, etc.)
		cartID := cart.CartID
		productID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		quantity := 1

		// Call the AddOrUpdateCartItem method
		status, err := s.addOrUpdateCartItemService(cartID, productID, quantity)
		if err != nil {
			// Handle the error (e.g., return an error response)
			http.Error(w, fmt.Sprintf("Failed to add/update cart item: %v", err), status)
			return
		}

		// Success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cart item added/updated successfully"))
	}
}
