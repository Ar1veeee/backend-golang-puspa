package middlewares

import (
	"backend-golang/shared/types"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimiterIP(rate time.Duration, capacity int64) gin.HandlerFunc {
	visitors := make(map[string]*ratelimit.Bucket)
	var mu sync.Mutex

	go func() {
		for {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			for ip, bucket := range visitors {
				if bucket.Available() == capacity {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()

		bucket, found := visitors[ip]
		if !found {
			bucket = ratelimit.NewBucket(rate, capacity)
			visitors[ip] = bucket
		}
		mu.Unlock()

		if bucket.TakeAvailable(1) == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, types.ErrorResponse{
				Success: false,
				Message: "Too Many Requests",
				Errors:  map[string]string{"errors": "You have made too many requests in a short period. Please try again later."},
			})
			return
		}

		c.Next()
	}
}
