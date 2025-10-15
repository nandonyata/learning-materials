package directives

import (
	"context"
	"fmt"

	"learn-graphql-go-gorm/middlewares"

	"github.com/99designs/gqlgen/graphql"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	// Check if user exists in context
	_, err := middlewares.ForContext(ctx)
	if err != nil {
		fmt.Println("error from directives")
		return nil, fmt.Errorf("access denied: %v", err)
	}

	// User is authenticated, proceed
	return next(ctx)
}
