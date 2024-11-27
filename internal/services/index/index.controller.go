package index

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"text/template"
	"time"

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
	r.HandleFunc("/", homePageHandler).Methods("GET")
	demoPageRoute := r.HandleFunc(fmt.Sprintf("/demo"), demoPageHandler).Methods("GET")
	log.Println(demoPageRoute.URL())
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
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
