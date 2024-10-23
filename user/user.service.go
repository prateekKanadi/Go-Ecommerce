package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const usersBasePath = "users"

// SetupRoutes :
func SetupUserRoutes(r *mux.Router, apiBasePath string) {
	r.HandleFunc(fmt.Sprintf("%s/%s", apiBasePath, usersBasePath), usersHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/{id}", apiBasePath, usersBasePath), userHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/resetPass", apiBasePath, usersBasePath), resetPassHandler).Methods("POST")
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
		resetCredHelper(w, r, updateUser, userID)
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

func resetPassHandler(w http.ResponseWriter, r *http.Request) {
	_, userID, res, err := userHandlerPrecheck(r)

	if err != nil {
		w.WriteHeader(res)
		return
	}

	resetCredHelper(w, r, updatePassword, userID)
	w.WriteHeader(http.StatusOK)
}

type helperFuncDef func(User) error

func resetCredHelper(w http.ResponseWriter, r *http.Request, inputFunc helperFuncDef, userID int) {
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

	err = inputFunc(updatedUser)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func userHandlerPrecheck(r *http.Request) (*User, int, int, error) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])

	if err != nil {
		return nil, 0, http.StatusNotFound, err
	}

	user, err := getUser(userID)

	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}
	if user == nil {
		return nil, 0, http.StatusNotFound, err
	}

	return user, userID, 0, nil
}
