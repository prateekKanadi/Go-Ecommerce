package user

import (
	"database/sql"
	"log"

	"github.com/ecommerce/database"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	// Convert the password to a byte array
	passwordBytes := []byte(password)

	// Generate the hashed password
	hashed, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password back to a string
	return string(hashed), nil
}

func getUserByEmail(email string) (*User, error) {
	row := database.DbConn.QueryRow(`SELECT 
	userId, 	
	email,
	password	
	FROM users
	WHERE email = ?`, email)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func getUser(userID int) (*User, error) {
	row := database.DbConn.QueryRow(`SELECT 
	userId, 	
	email,
	password	
	FROM users
	WHERE userId = ?`, userID)

	user := &User{}
	err := row.Scan(
		&user.UserID,
		&user.Email,
		&user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func removeUser(userID int) error {
	_, err := database.DbConn.Exec(`DELETE FROM users where userId = ?`, userID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func getUserList() ([]User, error) {
	results, err := database.DbConn.Query(`SELECT 
	userId, 	 
	email,
	password	
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
			&user.Password)

		users = append(users, user)
	}
	return users, nil
}

func updateEmail(user User) error {
	_, err := database.DbConn.Exec(`UPDATE users SET 		 				 
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

func updatePassword(user User) error {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	result, err := database.DbConn.Exec(`UPDATE users SET 		 				 
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

func RegisterUser(user User) (int, error) {
	return insertUser(user)
}
func LoginUser(user User) (int, error) {
	reqUser, err := getUserByEmail(user.Email)
	if err != nil {
		return 0, err
	}
	isCredMisMatch := bcrypt.CompareHashAndPassword([]byte(reqUser.Password), []byte(user.Password))
	if isCredMisMatch != nil {
		return 0, isCredMisMatch
	}
	return 0, nil
}
func insertUser(user User) (int, error) {
	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	result, err := database.DbConn.Exec(`INSERT INTO users  
	(email,
	password) VALUES (?, ?)`,
		user.Email,
		hashedPass)
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
