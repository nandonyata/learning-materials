package auth

import (
	"context"

	"learn-graphql-go-gorm/internal/pkg/jwt"
	"learn-graphql-go-gorm/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		// Allow unauthenticated access
		if header == "" {
			c.Next()
			return
		}

		// Validate JWT token
		username, err := jwt.ParseToken(header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		// Retrieve user from DB
		user := users.User{Username: username}
		id, err := users.GetUserIDByUsername(username)
		if err != nil {
			c.Next()
			return
		}

		user.ID = uint(id)

		// Store user in context
		ctx := context.WithValue(c.Request.Context(), userCtxKey, &user)
		c.Request = c.Request.WithContext(ctx)

		// Continue
		c.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)

	return raw
}
