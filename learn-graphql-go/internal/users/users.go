package users

import (
	"database/sql"
	"fmt"
	database "learn-graphql-go/internal/pkg/db/migrations/psql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() error {
	// Hash the password first
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Use correct PostgreSQL parameter placeholders
	stmt, err := database.Db.Prepare("INSERT INTO Users (Username, Password) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the insert
	_, err = stmt.Exec(user.Username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %v", err)
	}

	return nil
}

func (user *User) Authenticate() (bool, error) {
	stmt, err := database.Db.Prepare("SELECT Password FROM Users WHERE Username = $1")
	if err != nil {
		return false, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	var hashedPassword string
	err = stmt.QueryRow(user.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// User not found
			return false, nil
		}
		return false, fmt.Errorf("failed to query password: %v", err)
	}

	isValid := CheckPasswordHash(user.Password, hashedPassword)
	return isValid, nil
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = $1")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}
