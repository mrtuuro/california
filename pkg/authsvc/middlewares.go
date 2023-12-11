package authsvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type Middleware func(AuthService) AuthService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next AuthService) AuthService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   AuthService
	logger log.Logger
}

func (mw loggingMiddleware) Authenticate(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Authenticate",
			"took", time.Since(begin),
			"err", err)
	}(time.Now())
	return mw.next.Authenticate(ctx)
}
