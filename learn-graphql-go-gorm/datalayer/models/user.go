package models

import (
	"fmt"

	"learn-graphql-go-gorm/graph/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

// BeforeCreate hooks hash password before inserting to database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	return nil
}

// CheckPasswordHash compares a raw password with its hashed value
func (u *User) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) UserToGraphQLUser() *model.User {
	userGQL := model.User{
		ID:   fmt.Sprintf("%d", u.ID),
		Name: u.Name,
	}
	return &userGQL
}

func UserListToGraphQLUserList(users []*User) []*model.User {
	usersGQL := []*model.User{}
	for _, u := range users {
		usersGQL = append(usersGQL, u.UserToGraphQLUser())
	}
	return usersGQL
}
