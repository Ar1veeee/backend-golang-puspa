package middlewares

import (
	"backend-golang/shared/types"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiterIP(client *redis.Client, rate time.Duration, capacity int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:ip:%s", ip)

		ctx := context.Background()

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("Redis error for IP %s: %v", ip, err)
			c.Next()
			return
		}

		if count == 1 {
			if err := client.Expire(ctx, key, rate).Err(); err != nil {
				log.Printf("Failed to set expiration for key %s: %v", key, err)
				c.Next()
				return
			}
		}

		if count > capacity {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, types.ErrorResponse{
				Success: false,
				Message: "Too Many Requests",
				Errors:  map[string]string{"error": "You have made too many requests in a short period. Please try again later."},
			})
			return
		}

		c.Next()
	}
}

func RateLimiterUserID(client *redis.Client, rate time.Duration, capacity int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, exists := c.Get("userId")
		if !exists {
			c.Next()
			return
		}

		userId := key.(string)
		redisKey := fmt.Sprintf("ratelimit:userid:%s", userId)

		ctx := context.Background()

		count, err := client.Incr(ctx, redisKey).Result()
		if err != nil {
			log.Printf("Redis error for userID %s: %v", userId, err)
			c.Next()
			return
		}

		if count == 1 {
			if err := client.Expire(ctx, redisKey, rate).Err(); err != nil {
				log.Printf("Failed to set expiration for key %s: %v", redisKey, err)
				c.Next()
				return
			}
		}

		if count > capacity {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, types.ErrorResponse{
				Success: false,
				Message: "Too Many Requests",
				Errors:  map[string]string{"error": "You have made too many requests in a short period."},
			})
			return
		}

		c.Next()
	}
}
