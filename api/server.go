package api

import (
	"context"
	"fmt"
	"github.com/edison-moreland/gonduit/api/routes/auth"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func createRouter() *mux.Router {
	router := mux.NewRouter()

	// Routes go here
	auth.AddProfileRoutes(router)

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
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server? Reason %v", err.Error())
	}

	return nil
}
