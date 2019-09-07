package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Profile is a subset of the user model meant for presentation
type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// ProfileFromUser creates a new profile object from a user
func ProfileFromUser(u User, following bool) Profile {
	return Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}
}

// User is a database model to hold user information
type User struct {
	ID       uint   `json:"-" gorm:"unique;not null;primary_key"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Bio      string `json:"bio"`
	Image    string `json:"image"` // image url?
	Hash     string `json:"-" gorm:"not null"`
	Token    string `json:"token" gorm:"-"`

	// Self referencing many 2 many relationship
	// http://gorm.io/docs/many_to_many.html
	Following []User `json:"-" gorm:"many2many:followings;association_jointable_foreignkey:follower_id"`
}

// GetUser retrieves a user model from the database by it's username
func GetUser(username string) (User, error) {
	db := getDB()

	// Query db for user with specific username
	user := User{} // Possible SQL injection?
	db.Where(User{Username: username}).First(&user)

	if user.Username != username {
		// Query didn't return the correct user, probably doesn't exist
		return User{}, errors.New("user does not exist")
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

	// The actual error doesn't, we only care if it exists
	// If there is an error, the password isn't valid
	return err == nil
}

// UpdatePassword hashes a new password and updates the stored hash (DOES NOT SAVE TO DB)
func (u *User) UpdatePassword(newPassword string) error {
	// hash and salt password
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.MinCost)
	if err != nil {
		return fmt.Errorf("could not hash password. Reason: %v", err.Error())
	}

	// Turn it to a string and update field
	u.Hash = string(hash)
	println(u.Hash)

	return nil
}

// UpdateUser updates fields on the User model, does not update password. (Does not save to db)
func (u *User) UpdateUser(newUser User) {
	// TODO: There are libraries to merge structs, something like that should be used in the future

	// Manually merge each field we care about
	if newUser.Username != "" {
		u.Username = newUser.Username
	}
	if newUser.Email != "" {
		u.Email = newUser.Email
	}
	if newUser.Bio != "" {
		u.Bio = newUser.Bio
	}
	if newUser.Image != "" {
		u.Image = newUser.Image
	}
}

// GetProfile returns the profile for a user
func (u *User) GetProfile(following bool) Profile {
	return ProfileFromUser(*u, following)
}

// FollowUser adds a new user to this users following list
func (u *User) FollowUser(username string) error {
	userToFollow, err := GetUser(username)
	if err != nil {
		return err // User does't exist
	}

	getDB().Model(&u).Association("Following").Append(userToFollow)

	//
	// u.Following = append(u.Following, userToFollow)
	return nil
}

// UnFollowUser removes a user to this users following list
func (u *User) UnFollowUser(username string) error {
	userToFollow, err := GetUser(username)
	if err != nil {
		return err // User does't exist
	}

	getDB().Model(&u).Association("Following").Delete(userToFollow)

	return nil
}

// populateFollowing makes a database call to populate the "following" field of this model
func (u *User) populateFollowing() {
	getDB().Model(&u).Association("Following").Find(&u.Following)
}

// IsFollowingUser return bool indicating whether a user is in this users following list
func (u *User) IsFollowingUser(username string) bool {
	u.populateFollowing()

	for _, followee := range u.Following {
		if followee.Username == username {
			return true
		}
	}

	return false
}
