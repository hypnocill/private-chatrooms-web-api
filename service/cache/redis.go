package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func OpenRedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISTOGO_URL"),
		Password: "", // no password set (add password)
		DB:       0,  // use default DB
	})

	return rdb
}
