package session

import (
	"encoding/gob"
	"fmt"
	"math/rand"
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

	Cart struct {
		CartID    int
		Items     []CartItem
		CartTotal float64
	}

	CartItem struct {
		ID           int
		CartID       int
		ProductID    int
		Quantity     int
		ProductName  string
		PricePerUnit float64
		TotalPrice   float64
	}

	// Declare a custom type for the map
	IDCountMap map[int]int
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
	gob.Register(&Cart{})
	gob.Register(new(IDCountMap))
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

func InitAnonUserSession(sess *sessions.Session) {
	// Seed the random number generator with the current time
	anonUserId := rand.Int() // generates a random integer

	sess.Values["userId"] = anonUserId
	sess.Values["isAnon"] = true

	// Initialize `cart` before using it
	cart := &Cart{
		CartID: anonUserId,
	}

	// Initialize `user` before using it
	user := &User{
		UserID:   anonUserId,
		IsAdmin:  0,
		Email:    "",
		Password: "",
	}

	// Initialize the map using make (allocate memory for the map)
	IDCountMap := new(IDCountMap)

	sess.Values["IDCountMap"] = &IDCountMap
	sess.Values["user"] = &user
	sess.Values["cart"] = &cart
}
