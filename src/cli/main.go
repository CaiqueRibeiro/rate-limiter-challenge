package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/config"
	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/internal/infra/database"
)

func main() {
	token := flag.String("token", "", "A token to be set as custom rate limiter")
	maxReq := flag.Int64("maxreq", 0, "The max request token can make in a period of time")

	flag.Parse()
	if *token != "" {
		fmt.Printf("Saving rate limiter token \"%s\" to allow %d requests...\n", *token, *maxReq)

		cfg, err := config.Load(".")
		if err != nil {
			panic(err)
		}

		redisDB, err := database.NewRedisDatabase(*cfg)
		if err != nil {
			panic("cannot connect to Redis")
		}

		key := fmt.Sprintf("token_max_req:%s", *token)
		redisDB.Client.Set(context.Background(), key, *maxReq, 0)

		fmt.Println("Token registered.")
	}
}
