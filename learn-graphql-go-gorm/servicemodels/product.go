package servicemodels

type BaseProduct struct {
	Name        string
	Description string
	Quantity    int
}

type CreateProduct struct {
	BaseProduct
}

type UpdateProduct struct {
	ID uint
	BaseProduct
}
