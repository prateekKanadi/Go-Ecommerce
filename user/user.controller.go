package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	SetupProdUserRoutes(r, "prod")
}

func SetupProdUserRoutes(r *mux.Router, apiBasePath string) {
	r.HandleFunc(fmt.Sprintf("%s/%s", apiBasePath, usersBasePath), usersProdHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/{id}", apiBasePath, usersBasePath), userProdHandler)
	r.HandleFunc(fmt.Sprintf("%s/%s/resetPass", apiBasePath, usersBasePath), resetPassProdHandler).Methods("POST")
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
		removeUserService(userID)
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
