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
	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, cartBasePath)
	prodCartRouter := r.PathPrefix(prodUrlPath).Subrouter()

	prodCartRouter.HandleFunc("/{id}", AddToCartProdHandler(s)).Methods(http.MethodGet)
}

func AddToCartProdHandler(s *CartService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("cart-sess.Values[\"user\"]: ", sess.Values["user"])

		cart := sess.Values["cart"].(*session.Cart)
		cartID := cart.CartID
		productID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		quantity := 1
		log.Println("cart-productId", productID)
		log.Println("cart-cartID", cartID)

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
