package cart

import (
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

		//extracting isAnon flag from session
		isAnon := sess.Values["isAnon"].(bool)
		// Validate user
		user, ok := sess.Values["user"].(*session.User)
		if !ok || user == nil {
			http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusBadRequest)
			return
		}
		if isAnon {
			log.Println("cartpage Anon userID : ", user.UserID)
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

			if isAnon {
				var cartTotal float64
				var items []session.CartItem

				// iterate over items slice
				for i, item := range cart.Items {
					if (item.Quantity) > 0 {
						product, res, err := s.ProductService.GetProductService(item.ProductID)
						if err != nil {
							log.Println(err)
							http.Error(w, err.Error(), res)
							return
						}
						cart.Items[i].ProductName = product.ProductName
						cart.Items[i].PricePerUnit = product.PricePerUnit
						totalPrice := cart.Items[i].PricePerUnit * float64(item.Quantity)
						cart.Items[i].TotalPrice = totalPrice
						cartTotal += totalPrice
						items = append(items, cart.Items[i])
					}
				}
				cart.CartTotal = cartTotal
				cartListObj := cart
				cartListObj.Items = items
				cartList := cartListObj

				//updating values in session
				sess.Values["cart"] = &cart

				// saving session
				err = sess.Save(r, w)

				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = tmpl.Execute(w, map[string]interface{}{"CartItems": cartList.Items,
					"IsAdmin": user.IsAdmin, "CartTotal": cartList.CartTotal, "isAnon": isAnon})

				if err != nil {
					log.Println("Template execution error:", err)
					http.Error(w, "Error rendering cart list page", http.StatusInternalServerError)
					return
				}
			} else {
				cartListObj, res, err := s.getAllCartItemsService(cartID)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), res)
					return
				}
				cartList := cartListObj

				err = tmpl.Execute(w, map[string]interface{}{"CartItems": cartList.Items,
					"IsAdmin": user.IsAdmin, "CartTotal": cartList.CartTotal, "isAnon": isAnon})

				if err != nil {
					log.Println("Template execution error:", err)
					http.Error(w, "Error rendering cart list page", http.StatusInternalServerError)
					return
				}
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
		//extracting isAnon flag from session
		isAnon := sess.Values["isAnon"].(bool)

		// Validate user
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

		// Validate IDCountMap
		IDCountMap, ok := sess.Values["IDCountMap"].(*session.IDCountMap)
		if !ok || IDCountMap == nil {
			http.Error(w, `{"success": false, "error": "IDCountMap not found"}`, http.StatusBadRequest)
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
			var isFormQuantityNotNull = (r.FormValue("quantity") != "")
			if isFormQuantityNotNull {
				quantity, err = strconv.Atoi(r.FormValue("quantity"))
				if err != nil {
					log.Println("Invalid quantity:", err)
					http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusBadRequest)
					return
				}
			}

			if isAnon {
				_, exists := (*IDCountMap)[productID]
				if exists {
					// Value to search for
					targetValue := productID

					// Iterate over the items slice
					for i, item := range cart.Items {
						if item.ProductID == targetValue {
							if isFormQuantityNotNull {
								cart.Items[i].Quantity = quantity
							} else {
								cart.Items[i].Quantity += quantity
							}
							log.Println("Cart updated for Cart-Item: ", item.ID)
							break
						}
					}

					//updating values in session
					sess.Values["cart"] = &cart
				} else {
					//generate cartItemID
					cartItemID := len(cart.Items)

					//update IDCountMap
					(*IDCountMap)[productID] += 1

					//initialize Cartitem object
					item := session.CartItem{
						ID:           cartItemID,
						CartID:       cartID,
						ProductID:    productID,
						Quantity:     quantity,
						ProductName:  "",
						PricePerUnit: 0.0,
						TotalPrice:   0.0, // Quantity * PricePerUnit
					}
					cart.Items = append(cart.Items, item)

					//updating values in session
					sess.Values["cart"] = &cart
					sess.Values["IDCountMap"] = &IDCountMap
					fmt.Println("New Cart-Item added :", cartItemID)
				}
				// saving session
				err = sess.Save(r, w)

				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				// Call the AddOrUpdateCartItem method
				status, err := s.AddOrUpdateCartItemService(cartID, productID, quantity, isFormQuantityNotNull)
				if err != nil {
					// Handle the error (e.g., return an error response)
					log.Println("Error adding/updating cart item:", err)
					http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), status)
					return
				}
			}

			// Success response
			// Set response content type

			if isFormQuantityNotNull {
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

		//extracting isAnon flag from session
		isAnon := sess.Values["isAnon"].(bool)

		// Validate user
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

		// Validate IDCountMap
		IDCountMap, ok := sess.Values["IDCountMap"].(*session.IDCountMap)
		if !ok || IDCountMap == nil {
			http.Error(w, `{"success": false, "error": "IDCountMap not found"}`, http.StatusBadRequest)
			return
		}

		//Get cart ID from session
		cartID := cart.CartID

		switch r.Method {
		case http.MethodPost:
			// Get cart-item ID from URL
			cartItemID, err := strconv.Atoi(mux.Vars(r)["id"])
			if err != nil {
				log.Println("Invalid cartItem ID:", err)
				http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), http.StatusNotFound)
				return
			}

			err = r.ParseForm()
			if err != nil {
				log.Println(err)
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			if isAnon {
				// Value to search for
				targetValue := cartItemID

				// Variable to hold the key
				var keyFound int
				var found bool
				var deletedItem session.CartItem

				// Iterate over the items slice
				for i, item := range cart.Items {
					if item.ID == targetValue {
						cart.Items[i].Quantity = -1
						deletedItem = cart.Items[i]
						keyFound = cart.Items[i].ProductID
						found = true
						break
					}
				}

				// Check if a key was found
				if found {
					fmt.Printf("{Name: %s; Quantity: %d} Item deleted successfully \n", deletedItem.ProductName, deletedItem.Quantity)
					delete(*IDCountMap, keyFound)
				} else {
					fmt.Printf("Value %d not found in the map\n", targetValue)
				}

				//updating values in session
				sess.Values["cart"] = &cart
				sess.Values["IDCountMap"] = &IDCountMap

				// saving session
				err = sess.Save(r, w)

				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

			} else {
				// Call removeCartItem method
				status, err := s.removeCartItem(cartID, cartItemID)
				fmt.Println("Item with cartId : ", cartID, " and cartItemID : ", cartItemID, " was removed successfully")
				if err != nil {
					// Handle the error (e.g., return an error response)
					log.Println("Error deleting cart item:", err)
					http.Error(w, fmt.Sprintf(`{"success": false, "error": "%v"}`, err), status)
					return
				}
			}

			http.Redirect(w, r, "/prod/cart", http.StatusFound)
		case http.MethodOptions:
			return
		default:
			log.Printf("Entering in default")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}

}
