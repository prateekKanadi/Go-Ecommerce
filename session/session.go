package session

import (
	"encoding/gob"
	"net/http"

	"github.com/ecommerce/configuration"
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

func Init(config *configuration.Config) *sessions.CookieStore {
	registerTypes()
	store := sessions.NewCookieStore([]byte(config.Session.SessionKey), nil)
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

// Helper function to get session from request context
func GetSessionFromContext(r *http.Request) *sessions.Session {
	config := configuration.Conf
	session, _ := r.Context().Value(config.Session.SessionContextKey).(*sessions.Session)
	return session
}
