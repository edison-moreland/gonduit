package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User is a database model to hold user information
type User struct {
	Username string `json:"username" gorm:"unique;not null;primary_key"`
	Email    string `json:"email" gorm:"unique;not null"`
	Bio      string `json:"bio"`
	Image    string `json:"image"` // image url?
	Hash     string `json:"-" gorm:"not null"`
	Token    string `json:"token" gorm:"-"`
}

// GetUser retrieves a user model from the database by it's username
func GetUser(username string) (User, error) {
	db := getDB()

	// Query db for user with specific username
	user := User{} // Possible SQL injection?
	db.Where(User{Username: username}).First(&user)

	if (user == User{}) {
		// Query returned empty user, probgably doesn't exist
		return User{}, errors.New("User does not exist")
	}

	return user, nil
}

// Save migrates any changes to database
func (u *User) Save() {
	db := getDB()

	// Create new record if it doesn't already exist otherwise, save it
	if db.NewRecord(u) {
		db.Create(&u)
	} else {
		db.Save(&u)
	}
}

// CheckPassword compares a given password to the users hashed password
func (u User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// UpdatePassword hashes a new password and updates the stored hash (DOES NOT SAVCE TO DB)
func (u *User) UpdatePassword(newPassword string) error {
	// hash and salt password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("Could not hash password. Reason: %v", err.Error())
	}

	// Turn it to a string and update field
	u.Hash = string(hash)
	println(u.Hash)

	return nil
}
