package main

import (
	"fmt"

	"github.com/edison-moreland/gonduit/models"
)

func main() {
	println("Starting...")

	models.InitializeDB(":memory:")
	defer models.StopDB()
	println("Database initialized...")

	user1 := models.User{Username: "Bob joe", Email: "Jane@bo.com"}
	user1.UpdatePassword("Password1")
	user1.Save()
	println("Added user...")

	dbUser, _ := models.GetUser(user1.Username)
	fmt.Printf("%#v \n", dbUser)
}
