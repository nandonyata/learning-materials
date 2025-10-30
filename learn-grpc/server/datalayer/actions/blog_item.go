package actions

import (
	"context"
	"learn-grpc/server/datalayer/models"

	"gorm.io/gorm"
)

type BlogItemActionInterface interface {
	Save(ctx context.Context, blogItem *models.BlogItem) error
	SoftDelete(ctx context.Context, blogItem *models.BlogItem) error
	FetchById(ctx context.Context, id uint) (*models.BlogItem, error)
	FetchList(ctx context.Context) ([]*models.BlogItem, error)
}

type BlogItemAction struct {
	db *gorm.DB
}

func NewBlogItemAction(db *gorm.DB) BlogItemActionInterface {
	return &BlogItemAction{
		db: db,
	}
}

func (r *BlogItemAction) Save(ctx context.Context, blogItem *models.BlogItem) error {
	return r.db.WithContext(ctx).Save(blogItem).Error
}

func (r *BlogItemAction) SoftDelete(ctx context.Context, blogItem *models.BlogItem) error {
	return r.db.WithContext(ctx).Delete(blogItem).Error
}

func (r *BlogItemAction) FetchById(ctx context.Context, id uint) (*models.BlogItem, error) {
	var blogItem models.BlogItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&blogItem).Error
	return &blogItem, err
}

func (r *BlogItemAction) FetchList(ctx context.Context) ([]*models.BlogItem, error) {
	var blogItems []*models.BlogItem
	err := r.db.WithContext(ctx).Find(&blogItems).Error
	return blogItems, err
}
