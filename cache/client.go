package cache

import "github.com/redis/go-redis/v9"

func BuildRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}