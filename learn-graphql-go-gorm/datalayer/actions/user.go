package actions

import (
	"learn-graphql-go-gorm/datalayer/models"

	"gorm.io/gorm"
)

type UserActionInterface interface {
	Save(user *models.User) error
}

type UserAction struct {
	db *gorm.DB
}

func NewUserAction(db *gorm.DB) UserActionInterface {
	return &UserAction{
		db: db,
	}
}

func (r *UserAction) Save(user *models.User) error {
	return r.db.Save(user).Error
}
