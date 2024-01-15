package authsvc

import (
	"context"
	"errors"
	"strings"

	"california/pkg/repository"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	Authenticate(ctx context.Context) error
}

var (
	ErrAuthentication          = errors.New("authentication failed")
	ErrNoAuthTokenHeader       = errors.New("no auth token in the header")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken            = errors.New("invalid token")
)

type authService struct {
	signingKey string
	store      repository.Store
}

func NewAuthService(store repository.Store, signingKey string) AuthService {
	return &authService{
		signingKey: signingKey,
		store:      store,
	}
}

func (s *authService) Authenticate(ctx context.Context) error {
	_, err := isAuthenticated(ctx, s.signingKey)
	if err != nil {
		return err
	}
	return nil
}

func isAuthenticated(ctx context.Context, signingKey string) (context.Context, error) {
	// Extract the JWT token from the request header and validate it.
	tokenString := ctx.Value("Authorization").(string)
	if tokenString == "" {
		return nil, ErrNoAuthTokenHeader
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token algorithm is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, "email", claims["Email"])
		ctx = context.WithValue(ctx, "userId", claims["userId"])
		return ctx, nil
	}
	return nil, ErrInvalidToken
}
