package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxFailedAttempts = 5
	lockoutDuration   = 5 * time.Minute
)

type RateLimiterService interface {
	CheckLoginRateLimit(ctx context.Context, identifier string) error
	IncrementFailedAttempts(ctx context.Context, identifier string)
	ClearFailedAttempts(ctx context.Context, identifier string)
}

type rateLimiterService struct {
	redisClient *redis.Client
}

func NewRateLimiterService(redisClient *redis.Client) RateLimiterService {
	return &rateLimiterService{redisClient: redisClient}
}

func (s *rateLimiterService) CheckLoginRateLimit(ctx context.Context, identifier string) error {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)

	failedAttempts, err := s.redisClient.Get(ctx, redisKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if failedAttempts >= maxFailedAttempts {
		return errors.New("too many login attempts")
	}
	return nil
}

func (s *rateLimiterService) IncrementFailedAttempts(ctx context.Context, identifier string) {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)
	s.redisClient.Incr(ctx, redisKey)
	s.redisClient.Expire(ctx, redisKey, lockoutDuration)
}

func (s *rateLimiterService) ClearFailedAttempts(ctx context.Context, identifier string) {
	redisKey := fmt.Sprintf("login_attempts:%s", identifier)
	s.redisClient.Del(ctx, redisKey)
}
