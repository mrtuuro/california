package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"california/internal/config"
	"california/pkg/authsvc"
	"california/pkg/repository"
	"github.com/go-kit/kit/log"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	cfg := config.NewConfig()

	var svc authsvc.AuthService
	signingKey := os.Getenv("SECRET_KEY")
	c := context.WithValue(context.Background(), "foo", "bar")
	{
		store := repository.NewMongoStore(cfg)
		svc = authsvc.NewAuthService(store, signingKey)
		svc = authsvc.LoggingMiddleware(logger)(svc)
	}

	var h http.Handler
	{
		h = authsvc.MakeAuthHTTPHandler(c, svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.AuthHttpAddr)
		errs <- http.ListenAndServe(cfg.AuthHttpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
