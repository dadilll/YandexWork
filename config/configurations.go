package config

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
