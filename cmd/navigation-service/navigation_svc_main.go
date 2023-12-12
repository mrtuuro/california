package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"california/internal/config"
	"california/pkg/navigationsvc"
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

	var svc navigationsvc.NavigationService
	signingKey := os.Getenv("SECRET_KEY")
	c := context.WithValue(context.Background(), "foo", "bar")
	{
		store := repository.NewMongoStore(cfg)
		svc = navigationsvc.NewNavigationService(store)
		svc = navigationsvc.AuthMiddleware(signingKey)(svc)
		svc = navigationsvc.LoggingMiddleware(logger)(svc)
	}

	var h http.Handler
	{
		h = navigationsvc.MakeHTTPHandler(c, svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.NavigationHttpAddr)
		errs <- http.ListenAndServe(cfg.NavigationHttpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
