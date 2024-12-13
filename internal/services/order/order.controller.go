package order

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

const (
	prodBasePath = "prod"
	orderPath    = "order"
)

// SetupRoutes :
func SetupOrderRoutes(r *mux.Router, s *OrderService) {
	// -------------------------PROD----------------------
	orderUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, orderPath)
	orderCheckoutRouter := r.PathPrefix(orderUrlPath).Subrouter()
	orderCheckoutRouter.HandleFunc("/orderSummary", initiateOrderHandler(s))
	orderCheckoutRouter.HandleFunc("/orderHistory", orderHistoryHandler(s))
}

func initiateOrderHandler(s *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request came for initiate checkout")
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := sess.Values["user"].(*session.User)
		var userID = user.UserID
		err = r.ParseForm()
		if err != nil {
			err := errors.New("Error parsing form data")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tmpl, err := template.ParseFiles("template/orderSummary.html")

		switch r.Method {
		case http.MethodPost:

			fmt.Println("Request Entered in the POST method")

			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading order summary page", http.StatusInternalServerError)
				return
			}

			//
			orderTotalStr := r.FormValue("orderTotal")
			orderTotal, err := strconv.ParseFloat(orderTotalStr, 64)
			orderValueStr := r.FormValue("orderValue")
			orderValue, err := strconv.ParseFloat(orderValueStr, 64)
			deliveryMode := r.FormValue("deliveryMode")
			paymentMode := r.FormValue("paymentMode")

			if err != nil {
				log.Println(err)
				return
			}

			// Address Details Start
			houseNo := r.FormValue("HouseNo")
			landmark := r.FormValue("Landmark")
			city := r.FormValue("City")
			state := r.FormValue("State")
			pincode := r.FormValue("Pincode")
			contactNo := r.FormValue("Contact")

			shippingAddress := fmt.Sprintf("House No: %s , Landmark: %s , City: %s , State: %s , Pincode: %s , Phone Number: %s", houseNo, landmark, city, state, pincode, contactNo)
			// Address Details End

			// Get cart_items details --- START---
			// Validate cart
			cart, ok := sess.Values["cart"].(*session.Cart)
			if !ok || cart == nil {
				http.Error(w, `{"success": false, "error": "Cart not found"}`, http.StatusBadRequest)
				return
			}

			//Get cart ID from session
			cartID := cart.CartID
			Cart, err := s.CartService.Repo.GetAllCartItems(cartID)

			fmt.Println("Cart Items length is : ", len(Cart.Items))
			lenCartItems := len(Cart.Items)

			// Get cart_items details --- END ---

			// Order Creation in DB
			// --Validate that cart is NOT EMPTY--
			if lenCartItems > 0 {
				orderId, res, err := s.createOrderService(userID, deliveryMode, paymentMode, orderValue, orderTotal, shippingAddress)
				orderDetail, res, err := s.getOrderService(orderId)
				_, res, err = s.createOrderItemsService(orderId, userID, orderDetail)

				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), res)
					return
				}
				//

				//Redirect to order summary
				redirectURL := fmt.Sprintf("/%s/%s/orderSummary?orderId=%d", prodBasePath, orderPath, orderId)
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
				return
			}

			 htmlMessage := `
			 <!DOCTYPE html>
			 <html lang="en">
			 <head>
				 <meta charset="UTF-8">
				 <meta name="viewport" content="width=device-width, initial-scale=1.0">
				 <title>Redirecting...</title>
				 <style>
					 body {
						 font-family: Arial, sans-serif;
						 text-align: center;
						 margin-top: 50px;
					 }
				 </style>
				 <script>
					 setTimeout(function() {
						 window.location.href = "/prod/users/dashboard";
					 }, 5000); // Redirect after 5 seconds
				 </script>
			 </head>
			 <body>
				 <h1>Oops, your cart is empty!</h1>
				 <p>You will be redirected to the dashboard in 5 seconds...</p>
			 </body>
			 </html>
			 `
			 w.Header().Set("Content-Type", "text/html")
			 fmt.Fprint(w, htmlMessage)
			 log.Println("Rendered intermediate message for empty cart with delayed redirect")
			 return


		case http.MethodGet:

			fmt.Println("Request Entered in the GET method")
			// Retrieve orderId from query parameters
			orderIdStr := r.URL.Query().Get("orderId")
			if orderIdStr == "" {
				http.Error(w, "Missing orderId parameter", http.StatusBadRequest)
				return
			}

			orderId, err := strconv.Atoi(orderIdStr)
			order, err := s.Repo.GetOrdersAndOrderItemsByOrderID(orderId)

			if err != nil {
				log.Println(err)
				return
			}

			err = tmpl.Execute(w, map[string]interface{}{"Order": order})
			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering order summary page", http.StatusInternalServerError)
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

func orderHistoryHandler(s *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Order History Page Flow")

		sess, err := session.GetSessionFromContext(r)
		if sess == nil || err != nil {
			log.Println("Session not found:", err)
			http.Error(w, "Session not found", http.StatusUnauthorized)
			return
		}

		user, ok := sess.Values["user"].(*session.User)
		if !ok {
			log.Println("Failed to cast session value to User")
			http.Error(w, "Failed to retrieve user information", http.StatusUnauthorized)
			return
		}

		userID := user.UserID
		orders, err := s.Repo.GetAllOrdersAndOrderItemsByUserID(userID)

		switch r.Method {
		case http.MethodGet:
			tmpl, err := template.ParseFiles("template/orderHistory.html")
			if err != nil {
				log.Println("Error loading template:", err)
				http.Error(w, "Error loading order history page", http.StatusInternalServerError)
				return
			}

			data := map[string]interface{}{
				"Orders": orders,
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				log.Println("Error rendering template:", err)
				http.Error(w, "Error rendering order history page", http.StatusInternalServerError)
				return
			}

		case http.MethodOptions:
			w.WriteHeader(http.StatusOK)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
