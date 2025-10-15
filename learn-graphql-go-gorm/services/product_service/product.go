package product_service

import (
	"context"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/servicemodels"
)

type ProductServiceInterface interface {
	Create(ctx context.Context, user *models.User, req servicemodels.CreateProduct) (*models.Product, error)
	GetProductByID(ctx context.Context, id uint) (*models.Product, error)
	Update(ctx context.Context, req servicemodels.UpdateProduct) (*models.Product, error)
	GetProductList(ctx context.Context) ([]*models.Product, error)
}

type ProductService struct {
	productAction actions.ProductActionInterface
}

func NewService(
	productAction actions.ProductActionInterface,
) ProductServiceInterface {
	return &ProductService{
		productAction: productAction,
	}
}

// Create inserts a new product into the database
func (s *ProductService) Create(ctx context.Context, user *models.User, req servicemodels.CreateProduct) (*models.Product, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		UserID:      user.ID,
	}
	err := s.productAction.Save(ctx, &product)

	return &product, err
}

// GetProductWithID fetches product based on the ID
func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*models.Product, error) {
	product, err := s.productAction.FetchById(ctx, id)
	return product, err
}

// Update updates an existing product
func (s *ProductService) Update(ctx context.Context, req servicemodels.UpdateProduct) (*models.Product, error) {
	product, err := s.GetProductByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	product.Name = req.Name
	product.Description = req.Description
	product.Quantity = req.Quantity

	err = s.productAction.Save(ctx, product)

	return product, err
}

// GetProductList fetches all products
func (s *ProductService) GetProductList(ctx context.Context) ([]*models.Product, error) {
	products, err := s.productAction.FetchList(ctx)
	return products, err
}
