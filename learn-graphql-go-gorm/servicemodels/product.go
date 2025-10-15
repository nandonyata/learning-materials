package servicemodels

type CreateProduct struct {
	Name        string
	Description string
	Quantity    int
}

type UpdateProduct = CreateProduct
