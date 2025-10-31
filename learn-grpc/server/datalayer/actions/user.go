package actions

import (
	"context"
	"learn-grpc/server/datalayer/models"

	"gorm.io/gorm"
)

// UserActionInterface defines what operations we can do with users
// This is an INTERFACE (contract) - not the actual implementation
//
// WHY USE INTERFACE?
// 1. Easy to test (can create mock implementations)
// 2. Can swap implementations (e.g., switch from Postgres to MongoDB)
// 3. Clear contract of what operations are available
type UserActionInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

// UserAction is the ACTUAL implementation of UserActionInterface
// This uses GORM to interact with PostgreSQL database
type UserAction struct {
	db *gorm.DB // Database connection
}

// NewUserAction creates a new UserAction instance
// This is called when setting up the server
//
// USAGE:
//
//	db := database.DB
//	userAction := actions.NewUserAction(db)
func NewUserAction(db *gorm.DB) UserActionInterface {
	return &UserAction{
		db: db,
	}
}

// Create inserts a new user into the database
//
// WHAT IT DOES:
// SQL equivalent: INSERT INTO users (username, email, password) VALUES (?, ?, ?)
//
// PARAMETERS:
// - ctx: Request context (for cancellation, timeouts, tracing)
// - user: The user object to save (must have Username, Email, Password)
//
// AFTER SUCCESS:
// - user.ID will be populated with the new user's ID
// - user.CreatedAt will be set automatically by GORM
//
// ERRORS:
// - Unique constraint violation if email/username already exists
// - Database connection errors
func (r *UserAction) Create(ctx context.Context, user *models.User) error {
	// WithContext: Respects cancellation, timeouts from ctx
	// Create: Inserts new record and populates user.ID
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail finds a user by their email address
//
// WHAT IT DOES:
// SQL equivalent: SELECT * FROM users WHERE email = ? LIMIT 1
//
// RETURNS:
// - *models.User: The user if found
// - error: gorm.ErrRecordNotFound if user doesn't exist
//
// COMMON USAGE:
//
//	user, err := userAction.FindByEmail(ctx, "john@example.com")
//	if err == gorm.ErrRecordNotFound {
//	    // User doesn't exist - maybe show "invalid credentials"
//	}
func (r *UserAction) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// STEP BY STEP:
	// 1. WithContext(ctx): Use request context
	// 2. Where("email = ?", email): Add WHERE clause (? prevents SQL injection)
	// 3. First(&user): Get first matching record, populate user variable
	// 4. Error: Returns error if not found or database error
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	return &user, err
}

// FindByID finds a user by their ID
//
// WHAT IT DOES:
// SQL equivalent: SELECT * FROM users WHERE id = ? LIMIT 1
//
// USAGE EXAMPLE:
//
//	user, err := userAction.FindByID(ctx, 123)
//	if err != nil {
//	    return status.Errorf(codes.NotFound, "user not found")
//	}
//	fmt.Println(user.Username)
func (r *UserAction) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	// Same as FindByEmail, but searches by ID
	// Primary key lookup is very fast (uses index)
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error

	return &user, err
}

// ARCHITECTURE NOTE:
// This layer (Actions) handles ALL database operations
// Benefits:
// 1. Services don't need to know SQL/GORM
// 2. Easy to add caching, logging, metrics here
// 3. Can switch databases without changing services
// 4. Reusable across different services

// EXAMPLE DATA FLOW:
// Client Request → gRPC Service → Action → Database
//                      ↓              ↓
//                  (Business)    (Data Access)
//                  (Logic)        (SQL Queries)
