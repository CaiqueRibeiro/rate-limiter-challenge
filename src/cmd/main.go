package main

import (
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/config"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/web"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/web/handlers"
)

func main() {
	configs, err := config.Load(".")
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

	server := web.NewServer(
		configs.WebServerPort,
		handlers,
	)

	server.Run()
}
