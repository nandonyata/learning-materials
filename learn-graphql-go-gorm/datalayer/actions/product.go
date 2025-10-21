package actions

import (
	"context"
	"learn-graphql-go-gorm/datalayer/models"

	"gorm.io/gorm"
)

type ProductFilter = func(query *gorm.DB) *gorm.DB

type ProductActionInterface interface {
	Save(ctx context.Context, product *models.Product) error
	FetchById(ctx context.Context, id uint) (*models.Product, error)
	FetchList(ctx context.Context, filters ...ProductFilter) ([]*models.Product, error)
	Count(ctx context.Context, filters ...ProductFilter) (int64, error)
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

func (r *ProductAction) FetchList(ctx context.Context, filters ...ProductFilter) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.WithContext(ctx).
		Preload("User").
		Scopes(filters...).
		Find(&products).Error
	return products, err
}

func (r *ProductAction) Count(ctx context.Context, filters ...ProductFilter) (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Scopes(filters...).Count(&count).Error
	return count, err

}
