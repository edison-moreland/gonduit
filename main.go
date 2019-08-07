package main

import (
	"encoding/json"
	"fmt"

	"github.com/edison-moreland/gonduit/authentication/jwt"
	"github.com/edison-moreland/gonduit/models"
)

// TODO: Move to Viper config
const jwtSigningKey = "SupoerSecret"
const jwtTimeToLive = 24 // hours

func main() {
	println("Starting...")

	// Start DB
	models.InitializeDB(":memory:")
	defer models.StopDB()
	println("Database initialized...")

	user1 := models.User{Username: "Bob joe", Email: "Jane@bo.com"}
	user1.UpdatePassword("Password1")
	user1.Save()
	println("Added user...")

	userjwt, _ := jwt.Generate(&user1)

	jwtuser, _ := jwt.Validate(userjwt, []byte(jwtSigningKey))
	fmt.Printf("%#v \n", jwtuser)

	jwt.Revoke(userjwt)
	jwtuser, err := jwt.Validate(userjwt, []byte(jwtSigningKey))
	if err != nil {
		println(err.Error())
	}

	dbUser, _ := models.GetUser(user1.Username)
	fmt.Printf("%#v \n", dbUser)
	jsonUser, _ := json.Marshal(dbUser)
	fmt.Println(string(jsonUser))
}
