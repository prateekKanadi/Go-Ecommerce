package session

import (
	"encoding/gob"
	"os"

	"github.com/gorilla/sessions"
)

// variable declaration
var (
	// global
	Store *sessions.CookieStore

	// local

)

// type declaration
type (
	M      map[string]interface{}
	Person struct {
		FirstName string
		LastName  string
		Email     string
		Age       int
	}

	User struct {
		UserID   string
		Email    string
		Password string
	}
)

func Init() *sessions.CookieStore {
	registerTypes()
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   0,
		Secure:   false,
		HttpOnly: true,
	}

	Store = store
	return store
}

func registerTypes() {
	// gob.Register(&Person{})
	gob.Register(&User{})
	// gob.Register(&M{})
}
