package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ecommerce/cors"
)

const usersBasePath = "users"

// SetupRoutes :
func SetupRoutes(apiBasePath string) {
	handleUsers := http.HandlerFunc(usersHandler)
	handleUser := http.HandlerFunc(userHandler)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, usersBasePath), cors.Middleware(handleUsers))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, usersBasePath), cors.Middleware(handleUser))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userList, err := getUserList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		usersJson, err := json.Marshal(userList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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

		_, err = insertUser(newUser)
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

func userHandler(w http.ResponseWriter, r *http.Request) {
	urlBaseKeyWord := usersBasePath + "/"
	urlPathSegments := strings.Split(r.URL.Path, urlBaseKeyWord)
	userID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user, err := getUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
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
		var updatedUser User
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updatedUser.UserID != userID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// err = updateUser(updatedUser)
		// if err != nil {
		// 	log.Print(err)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		removeUser(userID)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
