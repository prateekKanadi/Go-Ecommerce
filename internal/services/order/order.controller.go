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
	prodBasePath     = "prod"
	orderSummaryPath = "orderSummary"
)

// SetupRoutes :
func SetupOrderRoutes(r *mux.Router, s *OrderService) {
	// -------------------------PROD----------------------
	orderUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, orderSummaryPath)
	orderCheckoutRouter := r.PathPrefix(orderUrlPath).Subrouter()
	orderCheckoutRouter.HandleFunc("", initiateOrderHandler(s))
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
		//get order id
		var userID = user.UserID
		err = r.ParseForm()
		if err != nil {
			err := errors.New("Error parsing form data")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orderTotalStr := r.FormValue("orderTotal")
		orderTotal, err := strconv.ParseFloat(orderTotalStr, 64)
		orderValueStr := r.FormValue("orderValue")
		orderValue, err := strconv.ParseFloat(orderValueStr, 64)
		deliveryMode := r.FormValue("deliveryMode")
		paymentMode := r.FormValue("paymentMode")
		orderId, res, err := s.createOrderService(userID, deliveryMode, paymentMode, orderValue, orderTotal)
		orderDetail, res, err := s.getOrderService(orderId)
		orderDetails, res, err := s.createOrderItemsService(orderId, userID, orderDetail)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		switch r.Method {
		case http.MethodPost:
			tmpl, err := template.ParseFiles("template/orderSummary.html")
			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading order summary page", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, map[string]interface{}{"Order": orderDetails})
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
