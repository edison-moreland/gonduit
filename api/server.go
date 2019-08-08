package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func createRouter() *mux.Router {
	router := mux.NewRouter()

	// Routes go here

	return router
}

// StartServer creates a new https server and starts serving content
func StartServer(address string) *http.Server {
	router := createRouter()
	handler := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, router.

	server := &http.Server{
		// Add logging and panic recovery
		Handler: handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, router)),
		
		Addr: address,

	}
}
