package usersvc

import (
	"context"
	"errors"

	"california/internal/helpers"
	"california/pkg/model"
	"california/pkg/repository"
)

type UserService interface {
	Register(context.Context, *model.User) error
	Login(ctx context.Context, email string, password string) (string, error)
}

type userService struct {
	store repository.Store
}

var (
	ErrInconsistentIDs           = errors.New("inconsistent IDs")
	ErrAlreadyExists             = errors.New("already exists")
	ErrNotFound                  = errors.New("not found")
	ErrAuthentication            = errors.New("authentication failed")
	ErrPasswordEmailDoesNotMatch = errors.New("password and email does not match")
	//ErrInternalDb                = errors.New("internal db error")
)

// Register TODO Add here to create a refresh token and return it to client.
func (s *userService) Register(ctx context.Context, user *model.User) error {
	exists, err := s.store.UserExists(ctx, user.Email, user.PhoneNumber)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	// Here we create a token for the user and return it to the client.
	// Later on client will have to use this token to send requests to the server.
	token, err := helpers.GenerateToken(user.Email, user.PhoneNumber)
	if err != nil {
		return err
	}
	user.RefreshToken = token

	// We need to hash the password before storing it in the database.
	hashedPass, err := helpers.HashRegisterPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPass
	err = s.store.InsertUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	// Get user by given email.
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return "", ErrNotFound
	}

	// We need to verify the given password with user's password.
	err = helpers.CompareLoginPasswordAndHash(password, user.Password)
	if err != nil {
		return "", ErrPasswordEmailDoesNotMatch
	}

	// Email and password matched, so we generate an access token and return it to the client.
	token, err := helpers.GenerateToken(user.Email, user.PhoneNumber)
	if err != nil {
		return "", err
	}

	// Use the created token to update the user's refresh token.
	user.RefreshToken = token
	return token, nil
}

func NewUserService(store repository.Store) UserService {
	return &userService{
		store: store,
	}
}
