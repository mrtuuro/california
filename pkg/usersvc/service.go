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
}

type userService struct {
	store repository.Store
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrInternalDb      = errors.New("internal db error")
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

	token, err := helpers.GenerateToken(user.Email, user.PhoneNumber)
	if err != nil {
		return err
	}
	user.RefreshToken = token

	err = s.store.InsertUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(store repository.Store) UserService {
	return &userService{
		store: store,
	}
}
