package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) getUserByEmail(email string) (*User, error) {
	row := repo.db.QueryRow(`SELECT 
	userId, 	
	email,
	password,
	isAdmin	
	FROM users
	WHERE email = ?`, email)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no user found with email %s: %v", email, err)
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) getUser(userID int) (*User, error) {
	row := repo.db.QueryRow(`SELECT 
	userId, 	
	email,
	password,
	isAdmin	
	FROM users
	WHERE userId = ?`, userID)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) removeUser(userID int) error {
	_, err := repo.db.Exec(`DELETE FROM users where userId = ?`, userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (repo *UserRepository) getAllUsers() ([]User, error) {
	results, err := repo.db.Query(`SELECT 
	userId, 	 
	email,
	password,
	isAdmin	
	FROM users`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	users := make([]User, 0)
	for results.Next() {
		var user User
		results.Scan(&user.UserID,
			&user.Email,
			&user.Password,
			&user.IsAdmin)

		users = append(users, user)
	}
	return users, nil
}

func (repo *UserRepository) updateEmail(user User) error {
	_, err := repo.db.Exec(`UPDATE users SET 		 				 
		email=?		
		WHERE userId=?`,
		user.Email,
		user.UserID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (repo *UserRepository) updatePassword(user User) error {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result, err := repo.db.Exec(`UPDATE users SET 		 				 
		password=?		
		WHERE userId=?`,
		hashedPass,
		user.UserID)
	if err != nil {
		result.RowsAffected()
		log.Println(err.Error())
		return err
	}
	return nil
}

func (repo *UserRepository) updateUser(user User) error {
	err_e := repo.updateEmail(user)
	// err_p := updatePassword(user)

	if err_e != nil {
		log.Println(err_e.Error())
		return err_e
	}

	return nil
}

func (repo *UserRepository) addUser(user User) (int, error) {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	result, err := repo.db.Exec(`INSERT INTO users  
	(email,
	password,
	name) VALUES (?, ?, ?)`,
		user.Email,
		hashedPass,
		user.Name)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

// functions for service layer outside user pkg

func (repo *UserRepository) RegisterUser(user User) (int, error) {
	return repo.addUser(user)
}

/* Login user authentication */
func (repo *UserRepository) LoginUser(user User) (int, error) {
	existingUser, err := repo.getUserByEmail(user.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Check if user was found (i.e., existingUser is not nil)
	if existingUser == nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	//compare existing-hashed-pass and request-pass
	isCredMisMatchError := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if isCredMisMatchError != nil {
		return http.StatusUnauthorized, errors.New("incorrect password")
	}
	return http.StatusOK, nil
}

func (repo *UserRepository) createCartForUser(userID int) (int, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// Create a new cart record for the user
	query := `INSERT INTO carts (user_id) VALUES (?)`
	result, err := repo.db.Exec(query, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create cart for user %d: %v", userID, err)
	}

	cartID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	log.Printf("Cart %d created for user %d successfully", cartID, userID)
	return int(cartID), nil
}

func (repo *UserRepository) GetCartForUser(userID int) (int, error) {
	// Query to get the cart ID for the given user ID
	query := `SELECT id FROM carts WHERE user_id = ?`
	var cartID int
	err := repo.db.QueryRow(query, userID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a specific error if no cart exists for the user
			return 0, fmt.Errorf("no cart found for user %d", userID)
		}
		log.Printf("Failed to get cart for user %d: %v", userID, err)
		return 0, err
	}

	log.Printf("Cart %d retrieved for user %d successfully", cartID, userID)
	return cartID, nil
}

// helper functions

func hashPassword(password string) (string, error) {
	// Convert the string password to a byte slice
	passwordBytes := []byte(password)

	// Generate the hashed password
	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password back to a string
	return string(hashed), nil
}

func (repo *UserRepository) GetAddressByUserId(userId int) (Address, error) {
	var address Address
	row := repo.db.QueryRow(`
        SELECT 
            userId, 
            houseNo, 
            landmark, 
            city, 
            state, 
            pincode, 
            phoneNumber 
        FROM address 
        WHERE userId = ?`, userId)
	err := row.Scan(
		&address.UserID,
		&address.HouseNo,
		&address.Landmark,
		&address.City,
		&address.State,
		&address.Pincode,
		&address.PhoneNumber,
	)

	if err == sql.ErrNoRows {
		return address, fmt.Errorf("no address found for userId %d: %v", userId, err)
	} else if err != nil {
		log.Println(err)
		return address, err
	}
	return address, nil
}

func (repo *UserRepository) UpdateAddressByUserId(userId int, address Address) error {
	_, err := repo.db.Exec(`
		UPDATE address 
		SET 
			houseNo = ?, 
			landmark = ?, 
			city = ?, 
			state = ?, 
			pincode = ?, 
			phoneNumber = ? 
		WHERE userId = ?`,
		address.HouseNo,
		address.Landmark,
		address.City,
		address.State,
		address.Pincode,
		address.PhoneNumber,
		userId)
	if err != nil {
		log.Println("Error updating address:", err)
		return fmt.Errorf("failed to update address for userId %d: %v", userId, err)
	}

	return nil
}

func (repo *UserRepository) AddAddressByUser(userId int, address Address) (int, error) {
	// Execute the SQL INSERT query to add a new address for the given userId.
	result, err := repo.db.Exec(`
		INSERT INTO address 
		(userId, houseNo, landmark, city, state, pincode, phoneNumber) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userId,
		address.HouseNo,
		address.Landmark,
		address.City,
		address.State,
		address.Pincode,
		address.PhoneNumber)

	// If there was an error executing the query, return the error.
	if err != nil {
		log.Println("Error adding address:", err)
		return 0, fmt.Errorf("failed to add address for userId %d: %v", userId, err)
	}

	// Get the last inserted ID to return the newly created address's ID.
	insertedId, err := result.LastInsertId()
	if err != nil {
		log.Println("Error fetching last insert ID:", err)
		return 0, fmt.Errorf("failed to fetch the last insert ID: %v", err)
	}

	// Return the newly inserted address ID (int).
	return int(insertedId), nil
}
