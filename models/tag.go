package models

import (
	"encoding/json"
)

// Tag describe article content
type Tag struct {
	ID   uint   `json:"-" gorm:"unique;not null;primary_key"`
	Name string `json:"name" gorm:"unique;not null"`
}

// MarshalJSON encodes the tag as a single string
func (t *Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(&t.Name)
}

// UnmarshalJSON decodes a tag from a single string
func (t *Tag) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Name)
}

// Save migrates any changes to database
func (t *Tag) Save() {
	db := getDB()

	// Create new record if it doesn't already exist otherwise, save it
	if db.NewRecord(t) {
		db.Create(&t)
	} else {
		db.Save(&t)
	}
}

//// GetTag retrieves a Tag model from the database by it's name
// func GetTag(name string) (Tag, error) {
//	db := getDB()
//
//	// Query db for tag with specific name
//	tag := Tag{}
//	db.Where(Tag{Name: name}).First(&tag)
//
//	if tag.Name != name {
//		// Query didn't return the correct tag, probably doesn't exist
//		return tag, errors.New("tag does not exist")
//	}
//
//	return tag, nil
//}

// GetTags retrieves all Tags from the database
func GetTags() []Tag {
	db := getDB()

	var tags []Tag
	db.Find(&tags)

	return tags
}
