package redis

import (
	"backend-golang/internal/infrastructure/config"
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	once        sync.Once
	initErr     error
)

func GetRedisClient() (*redis.Client, error) {
	once.Do(func() {
		redisHost := config.GetEnv("REDIS_HOST", "127.0.0.1")
		redisPort := config.GetEnv("REDIS_PORT", "6379")
		redisPassword := config.GetEnv("REDIS_PASSWORD", "")

		addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

		redisClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: redisPassword,
			DB:       0,
		})

		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			initErr = fmt.Errorf("failed to connect to Redis: %v", err)
			return
		}
	})

	return redisClient, initErr
}
