package user

import (
	"encoding/json"
	"log"
	"net/http"
)

func getAllUsersService() ([]byte, int) {
	userList, err := getAllUsers()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	usersJson, err := json.Marshal(userList)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return usersJson, http.StatusOK
}

func getUserService(userID int) (*User, int, error) {
	user, err := getUser(userID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if user == nil {
		return nil, http.StatusNotFound, err
	}

	return user, http.StatusOK, nil
}

func GetUserByEmailService(email string) (*User, error, int) {
	user, err := getUserByEmail(email)

	if err != nil {
		log.Print(err)
		return nil, err, http.StatusBadRequest
	}
	return user, nil, http.StatusOK
}

func addUserService(newUser User) (int, error) {
	_, err := addUser(newUser)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}

func updateUserService(updatedUser User) (int, error) {
	err := updateUser(updatedUser)

	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}

func updatePasswordService(updatedUser User) (int, error) {
	err := updatePassword(updatedUser)

	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, nil
	}
	return http.StatusOK, err
}

func removeUserService(userID int) error {
	return removeUser(userID)
}
