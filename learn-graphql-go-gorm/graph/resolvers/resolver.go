package resolvers

import "learn-graphql-go-gorm/services/user_service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService user_service.UserServiceInterface
}
