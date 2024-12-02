package index

import (
	"log"
	"net/http"
	"os/exec"
	"text/template"
	"time"

	"math/rand"

	"github.com/ecommerce/internal/core/session"
	"github.com/gorilla/mux"
)

func ServeIndexPage() {
	time.Sleep(1 * time.Second) // Wait a second for the server to start
	err := exec.Command("cmd", "/C", "start", "http://localhost:5000").Run()
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}

// SetupRoutes :
func SetupIndexRoutes(r *mux.Router) {
	r.HandleFunc("/", homePageHandler()).Methods("GET")
}

func homePageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sess, err := session.GetSessionFromContext(r)
		if sess == nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Seed the random number generator with the current time
		anonUserId := rand.Int() // generates a random integer

		sess.Values["userId"] = anonUserId
		sess.Values["isAnon"] = true

		// saving session
		err = sess.Save(r, w)

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userId := sess.Values["userId"]
		isAnon := sess.Values["isAnon"].(bool)
		if isAnon {
			log.Println("homepage Anon userID : ", userId)
		}

		// Parse the template file (adjust path if necessary)
		tmpl, err := template.ParseFiles("template/homePage.html")
		if err != nil {
			http.Error(w, "Error loading home page", http.StatusInternalServerError)
			log.Println("Template parsing error:", err)
			return
		}

		switch r.Method {
		case http.MethodGet:
			// Execute the template, sending data if needed (or nil if not)
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, "Error rendering home page", http.StatusInternalServerError)
				log.Println("Template execution error:", err)
			}
		case http.MethodOptions:
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func demoPageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file (adjust path if necessary)
	tmpl, err := template.ParseFiles("template/demoPage.html")
	if err != nil {
		http.Error(w, "Error loading demo page", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		// Execute the template, sending data if needed (or nil if not)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error rendering demo page", http.StatusInternalServerError)
			log.Println("Template execution error:", err)
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
