package actions

import (
	"context"
	"learn-graphql-go-gorm/datalayer/models"

	"gorm.io/gorm"
)

type ProductActionInterface interface {
	Save(ctx context.Context, product *models.Product) error
	FetchById(ctx context.Context, id uint) (*models.Product, error)
	FetchList(ctx context.Context) ([]*models.Product, error)
}

type ProductAction struct {
	db *gorm.DB
}

func NewProductAction(db *gorm.DB) ProductActionInterface {
	return &ProductAction{
		db: db,
	}
}

func (r *ProductAction) Save(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductAction) FetchById(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ?", id).
		First(&product).Error
	return &product, err
}

func (r *ProductAction) FetchList(ctx context.Context) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Preload("User").
		Find(&products).Error
	return products, err
}
