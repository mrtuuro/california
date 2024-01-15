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

func (mw loggingMiddleware) DeleteUser(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "DeleteUser",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.DeleteUser(ctx)
}

func (mw loggingMiddleware) Register(ctx context.Context, user *model.User) (insertedUser *model.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Register",
			"email", user.Email,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.Register(ctx, user)
}

func (mw loggingMiddleware) Login(ctx context.Context, email string, password string) (insertedUser *model.User, err error) {
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

func (mw loggingMiddleware) GetMe(ctx context.Context) (user *model.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetMe",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.GetMe(ctx)
}

func (mw loggingMiddleware) UpdateUserInfo(ctx context.Context, user *model.User) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "UpdateUserInfo",
			"took", time.Since(begin),
			"email", user.Email,
			"err", err)
	}(time.Now())
	return mw.next.UpdateUserInfo(ctx, user)
}

func (mw loggingMiddleware) UpdateVehicleInfo(ctx context.Context, vehicle *model.Vehicle) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "UpdateVehicleInfo",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.UpdateVehicleInfo(ctx, vehicle)
}

func (mw loggingMiddleware) ListAllUsers(ctx context.Context) (users []*model.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListAllUsers",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.ListAllUsers(ctx)
}

func (mw loggingMiddleware) SearchUsers(ctx context.Context, name string) (users []*model.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SearchUsers",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.SearchUsers(ctx, name)
}

type authMiddleware struct {
	next       UserService
	signingKey string
	c          context.Context
}

func (aw authMiddleware) Register(ctx context.Context, user *model.User) (insertedUser *model.User, err error) {
	return aw.next.Register(ctx, user)
}

func (aw authMiddleware) Login(ctx context.Context, email string, password string) (user *model.User, err error) {
	return aw.next.Login(ctx, email, password)
}

func (aw authMiddleware) VehicleRegister(ctx context.Context, vehicle *model.Vehicle) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.VehicleRegister(ctx, vehicle)

}

func (aw authMiddleware) GetMe(ctx context.Context) (user *model.User, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.GetMe(ctx)
}

func (aw authMiddleware) UpdateUserInfo(ctx context.Context, user *model.User) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.UpdateUserInfo(ctx, user)
}

func (aw authMiddleware) UpdateVehicleInfo(ctx context.Context, vehicle *model.Vehicle) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.UpdateVehicleInfo(ctx, vehicle)
}

func (aw authMiddleware) ListAllUsers(ctx context.Context) (users []*model.User, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.ListAllUsers(ctx)
}

func (aw authMiddleware) DeleteUser(ctx context.Context) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.DeleteUser(ctx)
}

func (aw authMiddleware) SearchUsers(ctx context.Context, name string) (users []*model.User, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.SearchUsers(ctx, name)
}

func AuthMiddleware(signingKey string) Middleware {
	return func(next UserService) UserService {
		return &authMiddleware{
			next:       next,
			signingKey: signingKey,
		}
	}
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
