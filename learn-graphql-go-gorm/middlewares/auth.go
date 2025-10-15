package middlewares

import (
	"context"
	"fmt"
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/pkg/jwt"
	"learn-graphql-go-gorm/services/user_service"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware(
	userService user_service.UserServiceInterface,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		// Allow unauthenticated access
		if header == "" {
			c.Next()
			return
		}

		// Validate JWT token
		userId, err := jwt.ParseToken(header)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		// Retrieve user from DB
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		user, err := userService.FetchByID(c, uint(userIdInt))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Store user in context
		ctx := context.WithValue(c.Request.Context(), userCtxKey, user)
		c.Request = c.Request.WithContext(ctx)

		// Continue
		c.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) (*models.User, error) {
	raw := ctx.Value(userCtxKey)

	if raw == nil {
		return nil, fmt.Errorf("user not found in context")
	}

	user := raw.(*models.User)

	return user, nil
}
