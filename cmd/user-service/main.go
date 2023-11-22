package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"california/internal/config"
	"california/pkg/repository"
	"california/pkg/usersvc"
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

	var svc usersvc.UserService
	{
		store := repository.NewMongoStore(cfg)
		svc = usersvc.NewUserService(store)
		svc = usersvc.LoggingMiddleware(logger)(svc)
	}

	var h http.Handler
	{
		h = usersvc.MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", cfg.HttpAddr)
		errs <- http.ListenAndServe(cfg.HttpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
