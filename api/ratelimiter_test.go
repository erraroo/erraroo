package api

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/erraroo/erraroo/logger"
	"github.com/stretchr/testify/assert"

	"gopkg.in/redis.v3"
)

const (
	perIntervalQuota = 2
	minDifference    = 1
	interval         = 10 * time.Second
)

type RedisRateLimiter struct {
	maxPerInterval int
	minDelta       int64
	interval       time.Duration
	client         *redis.Client
}

func (r *RedisRateLimiter) Check(key string) (bool, error) {
	tokens, err := r.getTokens(key)
	if err != nil {
		return false, err
	}

	lenTokens := len(tokens)
	exceededIntervalQuota := lenTokens >= r.maxPerInterval

	if lenTokens > 0 {
		lastRequestAt := tokens[len(tokens)-1]
		secondsSinceLastRequest := time.Now().UnixNano() - lastRequestAt
		if exceededIntervalQuota || secondsSinceLastRequest < r.minDelta {
			return false, nil
		}
	}

	return true, nil

}

func (r *RedisRateLimiter) getTokens(key string) ([]int64, error) {
	multi := r.client.Multi()
	defer multi.Close()

	now := time.Now()
	clearBefore := now.Add(-1 * r.interval).UnixNano()

	multi.ZRemRangeByScore(key, "0", fmt.Sprintf("%d", clearBefore))

	result, err := multi.ZRangeWithScores(key, 0, -1).Result()
	results := make([]int64, len(result))

	for i, r := range result {
		t, err := strconv.ParseInt(r.Member, 0, 64)
		if err != nil {
			return nil, err
		}

		results[i] = t
	}

	multi.ZAdd(key, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: fmt.Sprintf("%d", now.UnixNano()),
	})

	multi.Expire(key, interval)

	if err != nil {
		logger.Error("ratelimiting", "err", err)
		return nil, err
	}

	return results, nil

}

func TestRatelimiter(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	limiter := RedisRateLimiter{
		maxPerInterval: 1,
		minDelta:       0,
		interval:       10 * time.Millisecond,
		client:         client,
	}

	key := fmt.Sprintf("erraroo.test.ratelimting.%d", time.Now().UnixNano())

	ok, err := limiter.Check(key)
	assert.Nil(t, err)
	assert.Equal(t, ok, true)

	ok, err = limiter.Check(key)
	assert.Nil(t, err)
	assert.Equal(t, ok, false)

	time.Sleep(10 * time.Millisecond)

	ok, err = limiter.Check(key)
	assert.Nil(t, err)
	assert.Equal(t, ok, true)
}
