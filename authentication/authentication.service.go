package authentication

import (
	"encoding/json"
	"fmt"
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

		_, err = user.RegisterUser(newUser)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
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

		res, err := user.LoginUser(existingUser)
		if err != nil {
			log.Print(err)
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
