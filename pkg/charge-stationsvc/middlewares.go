package charge_stationsvc

import (
	"context"
	"strings"
	"time"

	"california/pkg/model"
	"california/pkg/usersvc"
	"github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware func(service StationService) StationService

type loggingMiddleware struct {
	next   StationService
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next StationService) StationService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw loggingMiddleware) StationRegister(ctx context.Context, station *model.Station) (insertedStation *model.Station, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "RegisterStation",
			"station_name", station.Name,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.StationRegister(ctx, station)
}

type authMiddleware struct {
	next       StationService
	signingKey string
	c          context.Context
}

func (aw authMiddleware) StationRegister(ctx context.Context, station *model.Station) (insertedStation *model.Station, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.StationRegister(ctx, station)
}

func AuthMiddleware(signingKey string) Middleware {
	return func(next StationService) StationService {
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
		return nil, usersvc.ErrNoAuthTokenHeader
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token algorithm is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, usersvc.ErrUnexpectedSigningMethod
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, usersvc.ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, "email", claims["Email"])
		ctx = context.WithValue(ctx, "userId", claims["userId"])
		return ctx, nil
	}
	return nil, usersvc.ErrInvalidToken
}
