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

	"github.com/ecommerce/internal/core/session"
	"github.com/ecommerce/utils"
	"github.com/gorilla/mux"
)

const (
	usersBasePath = "users"
	prodBasePath  = "prod"
	apiBasePath   = "api"
)

// SetupRoutes :
func SetupUserRoutes(r *mux.Router, s *UserService) {
	apiUrlPath := fmt.Sprintf("/%s/%s", apiBasePath, usersBasePath)
	userRouter := r.PathPrefix(apiUrlPath).Subrouter()

	userRouter.HandleFunc("", usersHandler(s))
	userRouter.HandleFunc("/{id}", userHandler(s))
	userRouter.HandleFunc("/resetPass", resetPassHandler(s)).Methods(http.MethodPost)

	// -------------------------PROD----------------------
	prodUrlPath := fmt.Sprintf("/%s/%s", prodBasePath, usersBasePath)
	prodUsersRouter := r.PathPrefix(prodUrlPath).Subrouter()

	// r.HandleFunc(fmt.Sprintf("/%s/%s", apiVersion, usersBasePath), usersProdHandler)
	// r.HandleFunc(fmt.Sprintf("/%s/%s/{id}", apiVersion, usersBasePath), userProdHandler)
	// r.HandleFunc(fmt.Sprintf("/%s/%s/resetPass", apiVersion, usersBasePath), resetPassProdHandler).Methods("POST")
	prodUsersRouter.HandleFunc("/dashboard", userDashboardHandler(s)).Methods(http.MethodGet)
}

func userDashboardHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if sess.Values == nil {
			err = errors.New("session values nil")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userId, err := session.GetSessionUserID(sess)
		if err != nil {
			log.Println("UserId is not set in session", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, res, err := s.getUserService(userId)

		if err != nil {
			log.Println("error : ", err)
			http.Error(w, err.Error(), res)
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
}

func usersProdHandler(w http.ResponseWriter, r *http.Request) {
}

func userProdHandler(w http.ResponseWriter, r *http.Request) {
}

func resetPassProdHandler(w http.ResponseWriter, r *http.Request) {
}

func usersHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userList, res, err := s.getAllUsersService()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}

			usersJson, err := json.Marshal(userList)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(usersJson)
			return
		case http.MethodPost:
			// add a new user to the list
			var newUser User
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(bodyBytes, &newUser)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if newUser.UserID != 0 {
				err := errors.New("UserId cannot be zero")
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			//adding user
			res, err := s.addUserService(newUser)

			if err == nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
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
}

func userHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		user, res, err := s.getUserService(userID)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		switch r.Method {
		case http.MethodGet:
			//return single user
			userJson, err := json.Marshal(user)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(userJson)
		case http.MethodPut:
			//update user in the list
			var updatedUser User
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(bodyBytes, &updatedUser)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if updatedUser.UserID != userID {
				err := errors.New("Payload User Id Mismatch")
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// update user cred
			res, err := s.updateUserService(updatedUser)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			res, err := s.removeUserService(userID)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), res)
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
}

func resetPassHandler(s *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		_, res, err := s.getUserService(userID)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}

		var updatedUser User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedUser)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if updatedUser.UserID != userID {
			err := errors.New("Payload User Id Mismatch")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//update user pass
		res, err = s.updatePasswordService(updatedUser)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), res)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
