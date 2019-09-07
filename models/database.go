package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Enabel sqlite dialect
)

// I don't know if this is a good pattern, but it's the best i've found
// DB connection object should be global to this package, but start/stoppable from the outside
// DB usage inside package should encapsulated so it can't be used with first initializing it
var _db *gorm.DB
var _dbIsInitialized = false

func getDB() *gorm.DB {
	// First verify database connection exists, and is still alive
	// Technically these are both recoverable, but I don't want to deal with it
	if !_dbIsInitialized {
		panic("Database isn't initialized!")
	}

	return _db
}

func migrateModels(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

// InitializeDB create a connection with the database and migrates all models
func InitializeDB(connection string) error {
	// ATM `connection` is just the sqlite path
	// in the future it will be a real db connection string

	if _dbIsInitialized {
		// Oops
		return errors.New("Database already initialized")
	}

	db, err := gorm.Open("sqlite3", connection)
	if err != nil {
		return errors.New("Can't initialize database")
	}

	db.LogMode(true)
	migrateModels(db)

	_db = db
	_dbIsInitialized = true

	return nil
}

// StopDB closes connection to database
func StopDB() {
	db := getDB()
	db.Close()
}
