package usersvc

import (
	"context"
	"errors"

	"california/internal/helpers"
	"california/pkg/model"
	"california/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	// Register and Login are public methods of the user.
	Register(context.Context, *model.User) error
	Login(ctx context.Context, email string, password string) (string, error)

	// VehicleRegister and VehicleUpdate are public methods of the vehicle.
	VehicleRegister(ctx context.Context, vehicle *model.Vehicle) error

	// GetMe is used to get the user's information and fill the blanks in the client.
	GetMe(ctx context.Context) (*model.User, error)

	// UpdateUserInfo is used to update the user's information.
	UpdateUserInfo(ctx context.Context, user *model.User) error

	// UpdateVehicleInfo is used to update the vehicle's information.
	UpdateVehicleInfo(ctx context.Context, vehicle *model.Vehicle) error
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
	ErrNoAuthToken               = errors.New("no auth token in the context")
	ErrInternalDb                = errors.New("internal db error")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrInvalidToken              = errors.New("invalid token")
)

// Register TODO Add here to create a refresh token and return it to client.
func (s *userService) Register(ctx context.Context, user *model.User) error {
	exists, err := s.store.UserExists(ctx, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	// Here we create a token for the user and return it to the client.
	// Later on client will have to use this token to send requests to the server.

	oidStr := primitive.NewObjectID().Hex()
	token, err := helpers.GenerateToken(user.Email, oidStr)
	if err != nil {
		return err
	}
	user.RefreshToken = token
	user.ID, _ = primitive.ObjectIDFromHex(oidStr)

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
	token, err := helpers.GenerateToken(user.Email, user.ID.Hex())
	if err != nil {
		return "", err
	}

	// Use the created token to update the user's refresh token.
	user.RefreshToken = token
	return token, nil
}

func (s *userService) VehicleRegister(ctx context.Context, vehicle *model.Vehicle) error {
	email := ctx.Value("email").(string)
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if err = s.store.InsertVehicleToUser(ctx, user, vehicle); err != nil {
		return err
	}
	return nil
}

func (s *userService) GetMe(ctx context.Context) (*model.User, error) {
	email := ctx.Value("email").(string)
	user, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, user *model.User) error {
	if err := s.store.UpdateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateVehicleInfo(ctx context.Context, vehicle *model.Vehicle) error {
	if err := s.store.UpdateVehicle(ctx, vehicle); err != nil {
		return err
	}
	return nil
}

func NewUserService(store repository.Store) UserService {
	return &userService{
		store: store,
	}
}
