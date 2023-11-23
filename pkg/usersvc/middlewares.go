package usersvc

import (
	"context"
	"strings"
	"time"

	"california/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(UserService) UserService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next UserService) UserService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   UserService
	logger log.Logger
}

func (mw loggingMiddleware) Register(ctx context.Context, user *model.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Register", "email", user.Email, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Register(ctx, user)
}

func (mw loggingMiddleware) Login(ctx context.Context, email string, password string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "Login", "email", email, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Login(ctx, email, password)
}

func (mw loggingMiddleware) VehicleRegister(ctx context.Context, vehicle *model.Vehicle) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "VehicleRegister",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.VehicleRegister(ctx, vehicle)
}

type authMiddleware struct {
	next       UserService
	signingKey string
	c          context.Context
}

func (aw authMiddleware) Register(ctx context.Context, user *model.User) (err error) {
	return aw.next.Register(ctx, user)
}

func (aw authMiddleware) Login(ctx context.Context, email string, password string) (token string, err error) {
	return aw.next.Login(ctx, email, password)
}

func (aw authMiddleware) VehicleRegister(ctx context.Context, vehicle *model.Vehicle) (err error) {
	// Extract the JWT token from the request header and validate it.
	tokenString := ctx.Value("Authorization").(string)
	if tokenString == "" {
		return ErrNoAuthToken
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token algorithm is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(aw.signingKey), nil
	})
	if err != nil {
		return ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, "email", claims["Email"])
		return aw.next.VehicleRegister(ctx, vehicle)
	}
	return ErrInvalidToken

}

func AuthMiddleware(signingKey string) Middleware {
	return func(next UserService) UserService {
		return &authMiddleware{
			next:       next,
			signingKey: signingKey,
		}
	}
}
