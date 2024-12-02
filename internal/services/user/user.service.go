package user

import (
	"log"
	"net/http"
)

// UserService handles business logic for user-related operations.
type UserService struct {
	Repo *UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) getAllUsersService() ([]User, int, error) {
	userList, err := s.Repo.getAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, http.StatusInternalServerError, err
	}
	return userList, http.StatusOK, nil
}

func (s *UserService) getUserService(userID int) (*User, int, error) {
	user, err := s.Repo.getUser(userID)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	if user == nil {
		return nil, http.StatusNotFound, err
	}

	return user, http.StatusOK, nil
}

func (s *UserService) GetUserByEmailService(email string) (*User, int, error) {
	user, err := s.Repo.getUserByEmail(email)

	if err != nil {
		log.Print(err)
		return nil, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

func (s *UserService) CreateCartForUserService(userID int) (int, int, error) {
	CartId, err := s.Repo.createCartForUser(userID)

	if err != nil {
		log.Print(err)
		return 0, http.StatusBadRequest, err
	}
	return CartId, http.StatusOK, nil
}

func (s *UserService) addUserService(newUser User) (int, error) {
	_, err := s.Repo.addUser(newUser)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *UserService) updateUserService(updatedUser User) (int, error) {
	err := s.Repo.updateUser(updatedUser)

	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *UserService) updatePasswordService(updatedUser User) (int, error) {
	err := s.Repo.updatePassword(updatedUser)

	if err != nil {
		log.Print(err)
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (s *UserService) removeUserService(userID int) (int, error) {
	err := s.Repo.removeUser(userID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (s *UserService) GetAddressByUserId(userID int) (Address, error) {
	address, err := s.Repo.GetAddressByUserId(userID)
	if err != nil {
		return address, err
	}
	return address, nil
}

func (s *UserService) UpdateAddressByUserId(userId int, address Address) error {
	err := s.Repo.UpdateAddressByUserId(userId, address)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AddAddressByUser(userId int, address Address) (int, error) {
	id, err := s.Repo.AddAddressByUser(userId, address)
	if err != nil {
		return id, err
	}
	return id, nil
}
