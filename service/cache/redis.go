package cache

import (
	"fmt"
	"os"
	"strconv"

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

	password := ""
	if passwordEnv := os.Getenv("REDIS_PASSWORD"); passwordEnv != "" {
		password = passwordEnv
	}

	database := "0"
	if databaseEnv := os.Getenv("REDIS_DATABASE"); databaseEnv != "" {
		database = databaseEnv
	}

	databaseStr, _ := strconv.Atoi(database)

	fmt.Println(host + ":" + port)
	fmt.Println(databaseStr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,    // no password set (add password)
		DB:       databaseStr, // use default DB
	})

	return rdb
}
