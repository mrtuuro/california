package usersvc

import (
	"context"
	"time"

	"california/pkg/model"
	"github.com/go-kit/kit/log"
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
