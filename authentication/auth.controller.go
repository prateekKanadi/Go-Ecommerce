package authentication

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ecommerce/user"
	"github.com/gorilla/mux"
)

const authBasePath = "auth"

// SetupRoutes :
func SetupAuthRoutes(r *mux.Router, apiBasePath string) {
	r.HandleFunc(fmt.Sprintf("%s/%s/login", apiBasePath, authBasePath), loginHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/register", apiBasePath, authBasePath), registerHandler)
	SetupProdAuthRoutes(r, "prod")
}

// SetupRoutes :
func SetupProdAuthRoutes(r *mux.Router, apiBasePath string) {
	r.HandleFunc(fmt.Sprintf("%s/%s/login", apiBasePath, authBasePath), loginProdHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/register", apiBasePath, authBasePath), registerProdHandler)
}

func registerProdHandler(w http.ResponseWriter, r *http.Request) {

}

func loginProdHandler(w http.ResponseWriter, r *http.Request) {

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// add a new user to the list
		var newUser user.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newUser.UserID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//register user
		res, err := registerUserService(newUser)

		if err == nil {
			w.WriteHeader(res)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Parse the template file (adjust path if necessary)
		tmpl, err := template.ParseFiles("template/login.html")
		if err != nil {
			http.Error(w, "Error loading login page", http.StatusInternalServerError)
			log.Println("Template parsing error:", err)
			return
		}

		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering login page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
	case http.MethodPost:
		// add a new product to the list
		var existingUser user.User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &existingUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if existingUser.UserID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//login user
		res, err := loginUserService(existingUser)

		if err == nil {
			w.WriteHeader(res)
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
