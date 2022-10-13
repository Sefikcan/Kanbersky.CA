package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/sefikcan/kanbersky.ca/pkg/config"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	redisHost := cfg.Redis.Url
	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})

	return client
}
