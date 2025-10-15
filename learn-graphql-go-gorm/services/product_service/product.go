package productservice

import (
	"context"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/graph/model"
	"learn-graphql-go-gorm/servicemodels"
)

type ProductServiceInterface interface {
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
func (s *ProductService) Create(ctx context.Context, user *models.User, req servicemodels.CreateProduct) (*model.Product, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		UserID:      user.ID,
	}
	err := s.productAction.Save(ctx, &product)

	return product.ProductToGraphQLProduct(), err
}
