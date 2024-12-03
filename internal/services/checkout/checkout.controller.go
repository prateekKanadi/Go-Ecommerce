package checkout

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

const (
	checkoutBasePath = "checkout"
	prodBasePath     = "prod"
	apiBasePath      = "api"
)

// SetupRoutes :
func SetupCheckoutRoutes(r *mux.Router, s *CheckoutService) {
	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, checkoutBasePath)
	prodCheckoutRouter := r.PathPrefix(prodUrlPath).Subrouter()
	prodCheckoutRouter.HandleFunc("", initiateCheckoutProdHandler(s))
}

func initiateCheckoutProdHandler(s *CheckoutService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request came for initiate checkout")
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := sess.Values["user"].(*session.User)

		// Fetch address details of the user
		var userId = user.UserID
		addressDetails, res, err := s.getAddressDetailsOfUser(userId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		switch r.Method {
		case http.MethodPost:
			tmpl, err := template.ParseFiles("template/checkout.html")
			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading product list page", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, map[string]interface{}{"Address": addressDetails})
			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering product list page", http.StatusInternalServerError)
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
