package state

import (
    "gopkg.in/redis.v5"
    "fmt"
)

var Redis *redis.Client

func NewRedisClient(url string) *redis.Client {
    redisOpts, err := redis.ParseURL(url)
    if err != nil {
        panic(fmt.Errorf("redis url parse error: %v", err))
    }
    return redis.NewClient(redisOpts)
}

// InitRedisClient redis client singleton
func InitRedisClient(url string) {
    Redis = NewRedisClient(url)
}
