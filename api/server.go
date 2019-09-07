package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/edison-moreland/gonduit/api/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func createRouter() *mux.Router {
	router := mux.NewRouter().PathPrefix("/api/").Subrouter()

	// Routes go here
	routes.AddUserRoutes(router)
	routes.AddProfileRoutes(router)

	return router
}

// StartServer creates a new https server and starts serving content
func StartServer(address string) *http.Server {
	router := createRouter()
	handler := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, router))

	server := &http.Server{
		// Add logging and panic recovery
		Handler: handler,

		Addr: address,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("HTTP server encoutered an error: %v", err.Error())
		}
	}()

	return server
}

// StopServer stops an https server with a timeout for existing connections
func StopServer(server *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server? Reason %v", err.Error())
	}

	return nil
}
