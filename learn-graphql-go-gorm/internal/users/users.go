package users

import (
	"errors"
	"fmt"
	database "learn-graphql-go-gorm/internal/pkg/db"

	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
}

// Create inserts a new user into the database with a hashed password
func (user *User) Create() error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = hashedPassword

	if err := database.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	return nil
}

// Authenticate checks if the username and password are valid
func (user *User) Authenticate() (bool, error) {
	var existing User
	err := database.DB.Where("username = ?", user.Username).First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user not found
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to query user: %v", err)
	}

	isValid := CheckPasswordHash(user.Password, existing.Password)
	return isValid, nil
}

// FetchByID fetches a user by ID
func (user *User) FetchByID(id uint) (*User, error) {
	var found User
	err := database.DB.First(&found, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by ID: %v", err)
	}
	return &found, nil
}

// GetUserIDByUsername returns the user ID for a given username
func GetUserIDByUsername(username string) (uint, error) {
	var user User
	err := database.DB.Select("id").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		log.Printf("error getting user by username: %v", err)
		return 0, err
	}
	return user.ID, nil
}

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a raw password with its hashed value
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
