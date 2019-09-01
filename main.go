package main

import (
	"github.com/edison-moreland/gonduit/api"
	"github.com/edison-moreland/gonduit/models"
	"os"
	"os/signal"
	"time"
)

// TODO: Move to Viper config
const apiAddress = "0.0.0.0:8080"
const apiExistingConnectionsTimeout = 5 * time.Second

func blockTillInterrupt() {
	// Capture interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Wait for SIGINT
	<-stop
}

func main() {
	println("Starting...")

	// Start DB
	err := models.InitializeDB(":memory:")
	defer models.StopDB()
	if err != nil {
		panic(err.Error())
	}
	println("Database initialized...")

	// Create test users
	user1 := models.User{Username: "Bob joe", Email: "Jane@bo.com"}
	_ = user1.UpdatePassword("Password1")
	user1.Save()

	user2 := models.User{Username: "Bofb joe", Email: "Janfe@bo.com"}
	_ = user2.UpdatePassword("Password1")
	_ = user2.FollowUser(user1.Username)
	user2.Save()

	println("Added user...")

	// Start https
	server := api.StartServer(apiAddress)
	defer func() {
		err := api.StopServer(server, apiExistingConnectionsTimeout)
		if err != nil {
			panic(err.Error())
		}
	}()
	println("Serving http...")
	println("API started!")

	// Now we wait
	blockTillInterrupt()
	println("SIGINT encountered, shutting down...")
}
