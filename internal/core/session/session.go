package session

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/ecommerce/configuration"
	"github.com/gorilla/sessions"
)

// type declaration
type (
	User struct {
		UserID   int
		Email    string
		Password string
		IsAdmin  int
	}
)

func Init(config *configuration.Config) (*sessions.CookieStore, error) {
	registerTypes()
	store := sessions.NewCookieStore([]byte(config.Session.SessionKey), nil)
	store.Options = &sessions.Options{
		Domain:   config.Session.Domain,
		Path:     config.Session.Path,
		MaxAge:   config.Session.MaxAge,
		Secure:   config.Session.Secure,
		HttpOnly: config.Session.HttpOnly,
	}

	return store, nil
}

func registerTypes() {
	gob.Register(&User{})
}

// Helper function to get session from request context
func GetSessionFromContext(r *http.Request) (*sessions.Session, error) {
	// Retrieve configuration from the context
	config, ok := r.Context().Value("config").(*configuration.Config)
	if !ok {
		return nil, fmt.Errorf("configuration not found in request context")
	}
	session, ok := r.Context().Value(config.Session.SessionContextKey).(*sessions.Session)
	if !ok {
		return nil, fmt.Errorf("session not found in request context")
	}
	return session, nil
}

// GetSessionUserID retrieves the userId from the session and returns an error if it doesn't exist or is not an int.
func GetSessionUserID(session *sessions.Session) (int, error) {
	userIdValue, exists := session.Values["userId"]
	if !exists {
		return 0, fmt.Errorf("userId not found in session")
	}

	userId, ok := userIdValue.(int)
	if !ok {
		return 0, fmt.Errorf("userId is not an int, found type: %T", userIdValue)
	}

	return userId, nil
}
