package product_service

import (
	"context"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/graph/model"
	"learn-graphql-go-gorm/servicemodels"

	"gorm.io/gorm"
)

type ProductServiceInterface interface {
	Create(ctx context.Context, user *models.User, req servicemodels.CreateProduct) (*models.Product, error)
	GetProductByID(ctx context.Context, id uint) (*models.Product, error)
	Update(ctx context.Context, req servicemodels.UpdateProduct) (*models.Product, error)
	GetProductList(ctx context.Context, query *servicemodels.ProductQuery) ([]*models.Product, error)
}

type ProductService struct {
	broadcaster   *EventBroadcaster
	productAction actions.ProductActionInterface
}

func NewService(
	productAction actions.ProductActionInterface,
) ProductServiceInterface {
	return &ProductService{
		broadcaster:   GetBroadcaster(),
		productAction: productAction,
	}
}

// Create inserts a new product into the database
func (s *ProductService) Create(ctx context.Context, user *models.User, req servicemodels.CreateProduct) (*models.Product, error) {
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    req.Quantity,
		ExpiredDate: req.ExpiredDate,
		UserID:      user.ID,
		User:        user,
	}
	err := s.productAction.Save(ctx, &product)

	if err == nil {
		s.broadcaster.BroadcastProductChange(model.ProductActionCreated, &product)
	}

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
	product.ExpiredDate = req.ExpiredDate

	err = s.productAction.Save(ctx, product)

	return product, err
}

// GetProductList fetches products based on query filter
func (s *ProductService) GetProductList(ctx context.Context, query *servicemodels.ProductQuery) ([]*models.Product, error) {
	var (
		products []*models.Product
		err      error
		filters  = []func(*gorm.DB) *gorm.DB{}
	)

	if query.IsPaginationApplicable() {
		if query.TotalData, err = s.productAction.Count(ctx, filters...); err != nil {
			return nil, err
		}

		filters = append(filters, servicemodels.WithPagination(query.PaginationRequest))
	}

	products, err = s.productAction.FetchList(ctx, filters...)
	return products, err
}
