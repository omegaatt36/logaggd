package cache

import (
	"log"

	"github.com/go-redis/redis/v7"
)

type redisClient struct {
	client *redis.Client
}

// Nil defines redis returned nil value error.
const Nil = redis.Nil

var singleton redisClient

// Initialize init package.
func Initialize(addr string) {
	singleton.client = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	err := Redis().Ping().Err()
	if err != nil {
		log.Panicf("connect to redis(%v) failed: %v", addr, err)
	}
}

// Redis returns redis client. It's safe to concurrent use.
func Redis() *redis.Client {
	if singleton.client == nil {
		panic("redis client is not created")
	}
	return singleton.client
}
