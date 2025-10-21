package servicemodels

import "time"

type BaseProduct struct {
	Name        string
	Description string
	Quantity    int
	ExpiredDate *time.Time
}

type CreateProduct struct {
	BaseProduct
}

type UpdateProduct struct {
	ID uint
	BaseProduct
}

type ProductQuery struct {
	PaginationRequest
}
