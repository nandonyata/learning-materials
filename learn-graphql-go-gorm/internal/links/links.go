package links

import (
	"fmt"
	database "learn-graphql-go-gorm/internal/pkg/db"
	"learn-graphql-go-gorm/internal/users"

	"log"
)

type Link struct {
	ID      uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title   string      `json:"title" gorm:"not null"`
	Address string      `json:"address" gorm:"not null"`
	UserID  uint        `json:"user_id" gorm:"not null"` // Foreign key
	User    *users.User `json:"user" gorm:"foreignKey:UserID"`
}

// Save creates a new Link record in the database
func (link *Link) Save() (uint, error) {
	if link.UserID == 0 && link.User.ID != 0 {
		link.UserID = link.User.ID
	}

	if link.UserID == 0 {
		return 0, fmt.Errorf("user ID is required to create a link")
	}

	if err := database.DB.Create(link).Error; err != nil {
		return 0, fmt.Errorf("failed to create link: %v", err)
	}

	log.Println("✅ Link inserted successfully!")
	return link.ID, nil
}

// GetAll retrieves all links from the database
func GetAll() ([]Link, error) {
	var links []Link
	if err := database.DB.Preload("User").Find(&links).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch links: %v", err)
	}
	return links, nil
}
