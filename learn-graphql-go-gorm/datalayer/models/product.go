package models

import (
	"fmt"
	"learn-graphql-go-gorm/graph/model"
	"learn-graphql-go-gorm/graph/scalars"
	"learn-graphql-go-gorm/servicemodels"
	"math"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(50);not null"`
	Description string     `gorm:"type:text"`
	Quantity    int        `gorm:"type:int"`
	ExpiredDate *time.Time `gorm:"type:date"`

	UserID uint  `gorm:"type:int"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (p *Product) ProductToGraphQLProduct() *model.Product {
	productGQL := model.Product{
		ID:          fmt.Sprintf("%d", p.ID),
		Name:        p.Name,
		Description: p.Description,
		Quantity:    int32(p.Quantity),
		ExpiredDate: scalars.NewDatePtr(p.ExpiredDate),
	}
	if p.User != nil {
		productGQL.User = p.User.UserToGraphQLUser()
	}
	return &productGQL
}

func ProductListToGraphQLProductList(products []*Product) []*model.Product {
	productsGQL := []*model.Product{}
	for _, p := range products {
		productsGQL = append(productsGQL, p.ProductToGraphQLProduct())
	}
	return productsGQL
}

func ProductListToPaginatedProducts(products []*Product, query servicemodels.ProductQuery) model.PaginatedProducts {
	response := model.PaginatedProducts{
		Products: ProductListToGraphQLProductList(products),
		Pagination: &model.PaginationInfo{
			TotalPages:  int32(math.Ceil(float64(query.TotalData) / float64(query.GetLimit()))),
			TotalData:   int32(query.TotalData),
			HasNextPage: int(query.TotalData) > (query.GetLimit() + query.GetOffset()),
		},
	}

	return response
}
