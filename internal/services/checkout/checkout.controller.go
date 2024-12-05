package checkout

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ecommerce/internal/core/session"
	User "github.com/ecommerce/internal/services/user"
	"github.com/gorilla/mux"
)

// type AddressData struct {
// 	HouseNo     string
// 	Landmark    string
// 	City        string
// 	State       string
// 	Pincode     string
// 	PhoneNumber string
// }

const (
	checkoutBasePath      = "checkout"
	prodBasePath          = "prod"
	apiBasePath           = "api"
	updateAddressBasePath = "/updateAddress"
)

// SetupRoutes :
func SetupCheckoutRoutes(r *mux.Router, s *CheckoutService) {
	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, checkoutBasePath)
	prodCheckoutRouter := r.PathPrefix(prodUrlPath).Subrouter()
	prodCheckoutRouter.HandleFunc("", initiateCheckoutProdHandler(s))
	prodCheckoutRouter.HandleFunc(updateAddressBasePath, updateAddressAtCheckoutProdHandler(s))
}


// Update Address at checkout Handler function
func updateAddressAtCheckoutProdHandler(s *CheckoutService) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Change address GET request being handled")
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("template/changeAddress.html")
		if err != nil {
			log.Println("Template parsing error:", err)
			http.Error(w, "Error loading product list page", http.StatusInternalServerError)
			return
		}
		user := sess.Values["user"].(*session.User)
		cart := sess.Values["cart"].(*session.Cart)
		_ = cart.CartID

		switch r.Method {
		case http.MethodGet:

			var userId = user.UserID
			addressDetails, res, err := s.getAddressDetailsOfUser(userId)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}
			err = tmpl.Execute(w, map[string]interface{}{"Address": addressDetails})

			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering product list page", http.StatusInternalServerError)
				return
			}
			return

		case http.MethodPost:
			return

		}
	}
}

// Initiate Checkout Handler
func initiateCheckoutProdHandler(s *CheckoutService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request came for initiate checkout")
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("template/checkout.html")
			if err != nil {
				log.Println("Template parsing error:", err)
				http.Error(w, "Error loading product list page", http.StatusInternalServerError)
				return
			}

		user := sess.Values["user"].(*session.User)
		cart := sess.Values["cart"].(*session.Cart)

		// Fetch address details of the user
		var userId = user.UserID
		addressDetails, res, err := s.getAddressDetailsOfUser(userId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		// Fetch Cart related data from cart table
		var cartId = cart.CartID
		cartData, err := s.getCartDetailsOfUser(cartId)
		fmt.Println("Details of cart fetched perfectly ...")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		fmt.Println("Cart Total is ", cartData.CartTotal)
		roundedoffTotal := fmt.Sprintf("%.2f", cartData.CartTotal)
		fmt.Println("Round off Total is ", roundedoffTotal)
		switch r.Method {
		case http.MethodGet:	
			err = tmpl.Execute(w, map[string]interface{}{"Address": addressDetails, "CheckoutData": roundedoffTotal})
			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering product list page", http.StatusInternalServerError)
				return
			}
			return

		case http.MethodPost:

			fmt.Println("Change address POST request being handled")
			// Take the updated address
			houseNo := r.FormValue("houseNo")
			landmark := r.FormValue("landmark")
			city := r.FormValue("city")
			state := r.FormValue("state")
			pincode := r.FormValue("pincode")
			phoneNumber := r.FormValue("phoneNumber")

			addressData := User.Address{
				HouseNo: houseNo,
				Landmark: landmark,
				City: city,
				State: state,
				Pincode: pincode,
				PhoneNumber: phoneNumber,
			}

			fmt.Println("All the data is received ",houseNo, landmark, city , pincode , state , phoneNumber)
			err = tmpl.Execute(w,map[string]interface{}{"Address":addressData, "CheckoutData": cartData.CartTotal})

			if err != nil {
				log.Println("Template execution error:", err)
				http.Error(w, "Error rendering checkout page", http.StatusInternalServerError)
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
