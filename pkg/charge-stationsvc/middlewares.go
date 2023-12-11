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

func (mw loggingMiddleware) FilterStation(ctx context.Context, brandName []string, socketType []string, currentType int) (stations []*model.Station, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "FilterStation",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.FilterStation(ctx, brandName, socketType, currentType)
}

func (mw loggingMiddleware) ListBrands(ctx context.Context) (brands []string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListBrands",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.ListBrands(ctx)
}

func (mw loggingMiddleware) ListSockets(ctx context.Context) (sockets []*model.Socket, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListSockets",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.ListSockets(ctx)
}

func (mw loggingMiddleware) StationRegister(ctx context.Context, station *model.Station) (insertedStation *model.Station, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "RegisterStation",
			"station_name", station.Brand,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.StationRegister(ctx, station)
}

func (mw loggingMiddleware) SearchStation(ctx context.Context, brandName string) (stations []*model.Station, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "SearchStation",
			"station_name", brandName,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.SearchStation(ctx, brandName)
}

func (mw loggingMiddleware) GetStations(ctx context.Context) (stations []*model.Station, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetAllStations",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.GetStations(ctx)
}

func (mw loggingMiddleware) UpdateStation(ctx context.Context, station *model.Station, stationId string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "UpdateStationInfo",
			"station_name", station.Brand,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.UpdateStation(ctx, station, stationId)
}

func (mw loggingMiddleware) RemoveStation(ctx context.Context, stationId string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "RemoveStation",
			"station_id", stationId,
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.RemoveStation(ctx, stationId)
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

func (aw authMiddleware) GetStations(ctx context.Context) (stations []*model.Station, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.GetStations(ctx)
}

func (aw authMiddleware) UpdateStation(ctx context.Context, station *model.Station, stationId string) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.UpdateStation(ctx, station, stationId)
}

func (aw authMiddleware) RemoveStation(ctx context.Context, stationId string) (err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return e
	}
	return aw.next.RemoveStation(ctx, stationId)
}

func (aw authMiddleware) SearchStation(ctx context.Context, brandName string) (stations []*model.Station, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.SearchStation(ctx, brandName)
}

func (aw authMiddleware) ListBrands(ctx context.Context) (brands []string, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.ListBrands(ctx)
}

func (aw authMiddleware) ListSockets(ctx context.Context) (sockets []*model.Socket, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.ListSockets(ctx)
}

func (aw authMiddleware) FilterStation(ctx context.Context, brandName []string, socketType []string, currentType int) (stations []*model.Station, err error) {
	ctx, e := isAuthenticated(ctx, aw.signingKey)
	if e != nil {
		return nil, e
	}
	return aw.next.FilterStation(ctx, brandName, socketType, currentType)
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
