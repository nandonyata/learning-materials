package servicemodels

type RegisterUser struct {
	Name     string
	Password string
}

type LoginUser = RegisterUser
