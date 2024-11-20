package middleware

import (
	"context"
	"net/http"

	"github.com/ecommerce/internal/core/setup"
	"github.com/gorilla/mux"
)

// CORS Middleware
func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS logic
			w.Header().Add("Access-Control-Allow-Origin", "*")
			// w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, r)
		})
	}
}

// AUTH Middleware
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS logic
			w.Header().Add("Access-Control-Allow-Origin", "*")
			// w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, r)
		})
	}
}

// Session Middleware : to load session and add it to the context
func SessionMiddleware(setupRes *setup.InitializationResult) func(http.Handler) http.Handler {
	//extract session store
	store := setupRes.Store
	config := setupRes.Config

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session logic
			session, err := store.Get(r, "session-name")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Add session to the request context
			ctx := context.WithValue(r.Context(), config.Session.SessionContextKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// InjectConfigMiddleware injects the configuration into the request context
func InjectConfigMiddleware(setupRes *setup.InitializationResult) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Inject configuration into context
			ctx := context.WithValue(r.Context(), "config", setupRes.Config)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Setup middleware
func RegisterMiddleWares(r *mux.Router, setupRes *setup.InitializationResult) {
	r.Use(InjectConfigMiddleware(setupRes)) // Add middleware for injecting config
	r.Use(CorsMiddleware())
	r.Use(SessionMiddleware(setupRes))
}
