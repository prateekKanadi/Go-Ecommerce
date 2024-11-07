package index

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// SetupRoutes :
func SetupIndexRoutes(r *mux.Router) {
	r.HandleFunc("/", homePageHandler).Methods("GET")
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