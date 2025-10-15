package actions

import (
	"context"
	"learn-graphql-go-gorm/datalayer/models"

	"gorm.io/gorm"
)

type UserActionInterface interface {
	Save(ctx context.Context, user *models.User) error
	FetchByName(ctx context.Context, username string) (*models.User, error)
	FetchById(ctx context.Context, id uint) (*models.User, error)
	FetchList(ctx context.Context) ([]*models.User, error)
}

type UserAction struct {
	db *gorm.DB
}

func NewUserAction(db *gorm.DB) UserActionInterface {
	return &UserAction{
		db: db,
	}
}

func (r *UserAction) Save(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserAction) FetchByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	return &user, err
}

func (r *UserAction) FetchById(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserAction) FetchList(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	return users, err
}
