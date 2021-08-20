package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func OpenRedisConnection() *redis.Client {
	host := "redis"
	if hostEnv := os.Getenv("REDIS_HOST"); hostEnv != "" {
		host = hostEnv
	}

	port := "6379"
	if portEnv := os.Getenv("REDIS_PORT"); portEnv != "" {
		port = portEnv
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "", // no password set (add password)
		DB:       0,  // use default DB
	})

	return rdb
}
