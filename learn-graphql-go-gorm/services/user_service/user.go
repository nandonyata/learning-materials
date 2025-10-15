package user_service

import (
	"context"
	"fmt"
	"learn-graphql-go-gorm/datalayer/actions"
	"learn-graphql-go-gorm/datalayer/models"
	"learn-graphql-go-gorm/graph/model"
	"learn-graphql-go-gorm/pkg/jwt"
	"learn-graphql-go-gorm/servicemodels"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req servicemodels.RegisterUser) (string, error)
	Login(ctx context.Context, req servicemodels.LoginUser) (string, error)
	FetchByID(ctx context.Context, id uint) (*model.User, error)
	FetchList(ctx context.Context) ([]*model.User, error)
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

// Register inserts a new user into the database with a hashed password
func (s *UserService) Register(ctx context.Context, req servicemodels.RegisterUser) (string, error) {
	user := models.User{
		Name:     req.Name,
		Password: req.Password,
	}

	if err := s.userAction.Save(ctx, &user); err != nil {
		return "", err
	}

	userIDStr := fmt.Sprintf("%d", user.ID)
	token, err := jwt.GenerateToken(userIDStr)
	return token, err
}

// Login fetch user and validate if the password is valid
func (s *UserService) Login(ctx context.Context, req servicemodels.LoginUser) (string, error) {
	user, err := s.userAction.FetchByName(ctx, req.Name)
	if err != nil {
		return "", nil
	}

	isValid := user.CheckPasswordHash(req.Password, user.Password)
	if isValid {
		return "", fmt.Errorf("invalid password")
	}

	userIDStr := fmt.Sprintf("%d", user.ID)
	token, err := jwt.GenerateToken(userIDStr)
	return token, err
}

// FetchByID fetches a user by ID
func (s *UserService) FetchByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userAction.FetchById(ctx, id)
	return s.userToGraphQLUser(user), err
}

// FetchList fetches list of users
func (s *UserService) FetchList(ctx context.Context) ([]*model.User, error) {
	users, err := s.userAction.FetchList(ctx)
	return s.userListToGraphQLUserList(users), err
}

func (s *UserService) userToGraphQLUser(user *models.User) *model.User {
	if user == nil {
		return nil
	}
	userGQL := model.User{
		ID:   fmt.Sprintf("%d", user.ID),
		Name: user.Name,
	}
	return &userGQL
}

func (s *UserService) userListToGraphQLUserList(users []*models.User) []*model.User {
	usersGQL := []*model.User{}
	for _, u := range users {
		usersGQL = append(usersGQL, s.userToGraphQLUser(u))
	}
	return usersGQL
}
