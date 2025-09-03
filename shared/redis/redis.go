package redis

import (
	"backend-golang/shared/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	redisHost := config.GetEnv("REDIS_HOST", "127.0.0.1")
	redisPort := config.GetEnv("REDIS_PORT", "6379")
	redisPassword := config.GetEnv("REDIS_PASSWORD", "")

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return client, nil
}
