package database

import (
	"context"
	"fmt"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabaseInterface interface{}

type RedisDatabase struct {
	Client *redis.Client
}

func NewRedisDatabase(
	cfg config.Conf,
) (*RedisDatabase, error) {
	addr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisDatabase{
		Client: client,
	}, nil
}
