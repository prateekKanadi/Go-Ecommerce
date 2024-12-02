package cart

import (
	"time"
	"encoding/json"
	"fmt"
	"html/template"
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

	prodCartRouter.HandleFunc("", cartsProdHandler(s))
	prodCartRouter.HandleFunc("/{id}", cartProdHandler(s))

	// Remove cart items Handler
	prodCartRouter.HandleFunc("/{id}/remove", removeCartItemProdHandler(s))
}

func cartsProdHandler(s *CartService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("cart-sess.Values[\"user\"]: ", sess.Values["user"])
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}

		// Validate cart
		cart, ok := sess.Values["cart"].(*session.Cart)
		if !ok || cart == nil {
			http.Error(w, `{"success": false, "error": "Cart not found"}`, http.StatusBadRequest)
			return
		}

		//Get cart ID from session
		cartID := cart.CartID

		switch r.Method {
		case http.MethodGet:
			tmpl, err := template.ParseFiles("template/cart_item_list.html")
			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading cart list page", http.StatusInternalServerError)
				return
			}

			cartList, res, err := s.getAllCartItemsService(cartID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}
			// cartList.CartTotal = 0.0
			err = tmpl.Execute(w, map[string]interface{}{"CartItems": cartList.Items,
				"IsAdmin": user.IsAdmin, "CartTotal": cartList.CartTotal})

			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering cart list page", http.StatusInternalServerError)
				return
			}
			return
		case http.MethodOptions:
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func cartProdHandler(s *CartService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("cart-sess.Values[\"user\"]: ", sess.Values["user"])
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}

		// Validate cart
		cart, ok := sess.Values["cart"].(*session.Cart)
		if !ok || cart == nil {
			http.Error(w, `{"success": false, "error": "Cart not found"}`, http.StatusBadRequest)
			return
		}

		//Get cart ID from session
		cartID := cart.CartID

		switch r.Method {
		case http.MethodPost:
			// Get product ID from URL
			productID, err := strconv.Atoi(mux.Vars(r)["id"])
			if err != nil {
				log.Println("Invalid product ID:", err)
				http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusNotFound)
				return
			}

			err = r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// extracting data of form values
			var quantity = 1
			if r.FormValue("quantity") != "" {
				quantity, err = strconv.Atoi(r.FormValue("quantity"))
				if err != nil {
					log.Println("Invalid quantity:", err)
					http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusBadRequest)
					return
				}
			}

			// hardcoded quantity set to 1
			// quantity := 1
			log.Println("cart-productId", productID)
			log.Println("cart-cartID", cartID)

			// Call the AddOrUpdateCartItem method
			status, err := s.addOrUpdateCartItemService(cartID, productID, quantity)
			if err != nil {
				// Handle the error (e.g., return an error response)
				log.Println("Error adding/updating cart item:", err)
				http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), status)
				return
			}

			// Success response
			// Set response content type

			if r.FormValue("quantity") != "" {
				http.Redirect(w, r, "/prod/cart", http.StatusSeeOther) // 303
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"success": true,
					"message": "Cart item added/updated successfully",
				})
			}

		case http.MethodOptions:
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}

/* remove cart items prod handler*/
func removeCartItemProdHandler(s *CartService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get session
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("cart-sess.Values[\"user\"]: ", sess.Values["user"])
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}

		// Validate cart
		cart, ok := sess.Values["cart"].(*session.Cart)
		if !ok || cart == nil {
			http.Error(w, `{"success": false, "error": "Cart not found"}`, http.StatusBadRequest)
			return
		}

		//Get cart ID from session
		cartID := cart.CartID

		
		switch r.Method {
		case http.MethodPost:
			// Get product ID from URL
			productID, err := strconv.Atoi(mux.Vars(r)["id"])
			if err != nil {
				log.Println("Invalid product ID:", err)
				http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusNotFound)
				return
			}

			err = r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			// extracting data of form values
			// var quantity = 1
			// if r.FormValue("quantity") != "" {
			// 	quantity, err = strconv.Atoi(r.FormValue("quantity"))
			// 	if err != nil {
			// 		log.Println("Invalid quantity:", err)
			// 		http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusBadRequest)
			// 		return
			// 	}
			// }

			// hardcoded quantity set to 1
			// quantity := 1
			log.Println("cart-productId", productID)
			log.Println("cart-cartID", cartID)

			// Call removeCartItem method
			status, err := s.removeCartItem(cartID, productID)
			fmt.Println("Item with cartId : ", cartID, " and cartId : ", productID, " was removed successfully")
			if err != nil {
				// Handle the error (e.g., return an error response)
				log.Println("Error deleting cart item:", err)
				http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), status)
				return
			}

			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode(map[string]interface{}{
			// 	"success": true,
			// 	"message": "Cart item deleted successfully"})
			time.Sleep(1 * time.Second)
			http.Redirect(w, r, "/prod/cart", http.StatusFound)

			// Call the AddOrUpdateCartItem method
			// status, err := s.addOrUpdateCartItemService(cartID, productID, quantity)
			// if err != nil {
			// 	// Handle the error (e.g., return an error response)
			// 	log.Println("Error adding/updating cart item:", err)
			// 	http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), status)
			// 	return
			// }

			// Success response
			// Set response content type

			// if r.FormValue("quantity") != "" {
			// 	http.Redirect(w, r, "/prod/cart", http.StatusSeeOther) // 303
			// } else {
			// 	w.Header().Set("Content-Type", "application/json")
			// 	json.NewEncoder(w).Encode(map[string]interface{}{
			// 		"success": true,
			// 		"message": "Cart item added/updated successfully",
			// 	})
			// }

		case http.MethodOptions:
			return
		default:
			log.Printf("Entering in default")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}
