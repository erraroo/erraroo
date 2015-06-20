package api

import (
	"fmt"
	"time"

	"github.com/erraroo/erraroo/logger"
	"gopkg.in/redis.v3"
)

type RateLimiter interface {
	Check(string, time.Duration, int) (bool, error)
}

var Limiter RateLimiter

type RedisRateLimiter struct {
	client *redis.Client
}

func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

func (r *RedisRateLimiter) Check(key string, interval time.Duration, max int) (bool, error) {
	tokens, err := r.getTokens(key, interval)
	if err != nil {
		return false, err
	}

	if len(tokens) >= max {
		return false, nil
	}

	return true, nil

}

func (r *RedisRateLimiter) getTokens(key string, interval time.Duration) ([]int64, error) {
	multi := r.client.Multi()
	defer multi.Close()

	now := time.Now()
	nano := now.UnixNano()

	clearBefore := now.Add(-1 * interval).UnixNano()
	multi.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", clearBefore))

	result, err := multi.ZRangeWithScores(key, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	results := make([]int64, len(result))
	for i, r := range result {
		results[i] = int64(r.Score)
	}

	multi.ZAdd(key, redis.Z{
		Score:  float64(nano),
		Member: fmt.Sprintf("%d", nano),
	})

	if interval > time.Second {
		multi.Expire(key, interval)
	} else {
		multi.Expire(key, 1*time.Second)
	}

	if err != nil {
		logger.Error("ratelimiting", "key", key, "err", err)
		return nil, err
	}

	return results, nil
}

func NoLimiter() RateLimiter {
	return &noLimiter{}
}

type noLimiter struct{}

func (*noLimiter) Check(key string, i time.Duration, n int) (bool, error) {
	return true, nil
}
