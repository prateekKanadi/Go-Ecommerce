package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/ecommerce/session"
	s "github.com/ecommerce/session"
	"github.com/ecommerce/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	usersBasePath = "users"
	apiVersion    = "prod"
	apiBasePath   = "api"
)

var (
	store *sessions.CookieStore
)

// SetupRoutes :
func SetupUserRoutes(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s/%s", apiBasePath, usersBasePath), usersHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/{id}", apiBasePath, usersBasePath), userHandler)
	r.HandleFunc(fmt.Sprintf("/%s/%s/resetPass", apiBasePath, usersBasePath), resetPassHandler).Methods("POST")

	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("%s/%s", apiVersion, usersBasePath)
	// r.HandleFunc(fmt.Sprintf("/%s/%s", apiVersion, usersBasePath), usersProdHandler)
	// r.HandleFunc(fmt.Sprintf("/%s/%s/{id}", apiVersion, usersBasePath), userProdHandler)
	// r.HandleFunc(fmt.Sprintf("/%s/%s/resetPass", apiVersion, usersBasePath), resetPassProdHandler).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/dashboard", prodUrlPath), userDashboardHandler).Methods("GET")
}

func userDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session := session.GetSessionFromContext(r)
	if session == nil {
		http.Error(w, errors.New("session not found in context").Error(), http.StatusInternalServerError)
		return
	}

	if session.Values == nil {
		http.Error(w, errors.New("session values nil").Error(), http.StatusInternalServerError)
		return
	}

	userId, err := s.GetSessionUserID(session)
	if err != nil {
		log.Println("UserId is not set in session")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, res, err := getUserService(userId)

	if err != nil {
		http.Error(w, err.Error(), res)
		log.Println("error : ", err)
		return
	}
	log.Println(utils.ToString(*user))

	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/dashboard.html")
	if err != nil {
		http.Error(w, "Error loading dashboard page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, user)
		if err != nil {
			http.Error(w, "Error rendering dashboard page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
			return
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func usersProdHandler(w http.ResponseWriter, r *http.Request) {
}

func userProdHandler(w http.ResponseWriter, r *http.Request) {
}

func resetPassProdHandler(w http.ResponseWriter, r *http.Request) {
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usersJson, res := getAllUsersService()
		if usersJson == nil {
			w.WriteHeader(res)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(usersJson)
	case http.MethodPost:
		// add a new user to the list
		var newUser User
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

		//adding user
		res, err := addUserService(newUser)

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

func userHandler(w http.ResponseWriter, r *http.Request) {
	user, userID, res, err := userHandlerPrecheck(r)

	if err != nil {
		w.WriteHeader(res)
		return
	}

	switch r.Method {
	case http.MethodGet:
		//return single user
		userJson, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
	case http.MethodPut:
		//update user in the list
		updatedUser, res := resetCredHelper(r, userID)

		if res != http.StatusOK {
			w.WriteHeader(res)
			return
		}

		//update user cred
		res, err := updateUserService(updatedUser)

		if err == nil {
			w.WriteHeader(res)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		err := removeUserService(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func resetPassHandler(w http.ResponseWriter, r *http.Request) {
	_, userID, res, err := userHandlerPrecheck(r)

	if err != nil {
		w.WriteHeader(res)
		return
	}

	updatedUser, res := resetCredHelper(r, userID)

	if res != http.StatusOK {
		w.WriteHeader(res)
		return
	}

	//update user pass
	res, err = updatePasswordService(updatedUser)

	if err == nil {
		w.WriteHeader(res)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func resetCredHelper(r *http.Request, userID int) (User, int) {
	var updatedUser User
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return updatedUser, http.StatusBadRequest
	}

	err = json.Unmarshal(bodyBytes, &updatedUser)
	if err != nil {
		return updatedUser, http.StatusBadRequest
	}

	if updatedUser.UserID != userID {
		return updatedUser, http.StatusBadRequest
	}

	return updatedUser, http.StatusOK
}

func userHandlerPrecheck(r *http.Request) (*User, int, int, error) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])

	if err != nil {
		return nil, 0, http.StatusNotFound, err
	}

	user, res, err := getUserService(userID)

	if err != nil {
		return nil, 0, res, err
	}
	if user == nil {
		return nil, 0, res, err
	}

	return user, userID, 0, nil
}
