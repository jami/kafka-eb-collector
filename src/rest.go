package src

import (
	"net/http"

	"github.com/gorilla/mux"
)

var router *mux.Router

func handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORSHeader(w)
		if r.Method == "OPTIONS" {
			return
		}

		res := []string{"create done"}

		HTTPAsJSON(w, res)
	}
}

// setupCORSHeader allow cors
func setupCORSHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-WP-Nonce")
}

// registerRouter creates the app router
func registerRouter() {
	router = mux.NewRouter()
	router.PathPrefix("/create").HandlerFunc(handleCreate())
	router.PathPrefix("/append").HandlerFunc(handleCreate())
}

// RestListenerListen provides rest api endpoint
func RestListenerListen(endpoint string) error {
	return http.ListenAndServe(endpoint, router)
}
