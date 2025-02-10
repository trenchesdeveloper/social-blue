package cache

import (
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	client  *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisConfig{
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisConfig{
		client: redisClient,
	}
}

func (r *RedisConfig) Set(key string, value interface{})  error {
	return nil
}

func (r *RedisConfig) Get(key string) (string, error) {
	return "", nil
}