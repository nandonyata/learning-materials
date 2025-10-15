package user_service

import (
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/datalayer/models"
)

type UserServiceInterface interface {
	FetchByID(id uint) (*models.User, error)
}

type UserService struct {
	userAction actions.UserActionInterface
}

func NewService(
	userAction actions.UserActionInterface,
) UserServiceInterface {
	return &UserService{
		userAction: userAction,
	}
}

// Create inserts a new user into the database with a hashed password
func (user *UserService) Create() error {
	// hashedPassword, err := HashPassword(user.Password)
	// if err != nil {
	// 	return fmt.Errorf("failed to hash password: %v", err)
	// }
	// user.Password = hashedPassword

	// if err := database.DB.Create(user).Error; err != nil {
	// 	return fmt.Errorf("failed to create user: %v", err)
	// }
	return nil
}

// Authenticate checks if the username and password are valid
func (user *UserService) Authenticate() (bool, error) {
	// var existing User
	// err := database.DB.Where("username = ?", user.Username).First(&existing).Error

	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	// user not found
	// 	return false, nil
	// }
	// if err != nil {
	// 	return false, fmt.Errorf("failed to query user: %v", err)
	// }

	// isValid := CheckPasswordHash(user.Password, existing.Password)
	// return isValid, nil
	return true, nil
}

// FetchByID fetches a user by ID
func (user *UserService) FetchByID(id uint) (*models.User, error) {
	// var found User
	// err := database.DB.First(&found, id).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to fetch user by ID: %v", err)
	// }
	// return &found, nil
	return &models.User{}, nil
}

// GetUserIDByUsername returns the user ID for a given username
func (userr *UserService) GetUserIDByUsername(username string) (uint, error) {
	// var user User
	// err := database.DB.Select("id").Where("username = ?", username).First(&user).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return 0, nil
	// }
	// if err != nil {
	// 	log.Printf("error getting user by username: %v", err)
	// 	return 0, err
	// }
	// return user.ID, nil
	return 1, nil
}
