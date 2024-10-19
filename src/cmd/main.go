package main

import (
	"time"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/config"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/database"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/web"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/web/handlers"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/web/middlewares"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		panic(err)
	}

	exampleHandler := handlers.NewExampleHandler()
	handlers := []web.Handler{
		{
			Path:        "/",
			Method:      "GET",
			HandlerFunc: exampleHandler.Get,
		},
	}

	redisDB, err := database.NewRedisDatabase(*cfg)
	if err != nil {
		panic("cannot connect to Redis")
	}

	redisStrategy := strategies.NewRedisLimiter(redisDB.Client, time.Now)
	rateLimiter := ratelimiter.NewRateLimiter(redisStrategy, cfg.IPMaxRequests, cfg.TokenMaxRequests, cfg.TimeWindowMilliseconds)
	rlMiddleware := middlewares.NewRateLimiterMiddleware(rateLimiter)

	middlewares := []web.Middleware{
		{
			Name:    "RateLimiter",
			Handler: rlMiddleware.Handle,
		},
	}

	server := web.NewServer(
		cfg.WebServerPort,
		handlers,
		middlewares,
	)

	server.Run()
}
