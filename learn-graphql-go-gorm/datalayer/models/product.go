package models

import (
	"fmt"
	"learn-graphql-go-gorm/graph/model"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);not null"`
	Description string `gorm:"type:text"`
	Quantity    int    `gorm:"type:int"`

	UserID uint  `gorm:"type:int"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (p *Product) ProductToGraphQLProduct() *model.Product {
	productGQL := model.Product{
		ID:          fmt.Sprintf("%d", p.ID),
		Name:        p.Name,
		Description: p.Description,
		Quantity:    int32(p.Quantity),
		User:        p.User.UserToGraphQLUser(),
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
