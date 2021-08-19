package cache

import (
	"github.com/go-redis/redis/v8"
)

func OpenRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set (add password)
		DB:       0,  // use default DB
	})

	return rdb
}
