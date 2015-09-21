package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gopkg.in/redis.v3"
)

func TestRatelimiter(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer client.Close()

	limiter := RedisRateLimiter{
		client: client,
	}

	key := fmt.Sprintf("erraroo.test.ratelimting.%d", time.Now().UnixNano())
	interval := 5 * time.Millisecond
	max := 1

	ok, err := limiter.Check(key, interval, max)
	assert.Nil(t, err)
	assert.Equal(t, ok, true)

	ok, err = limiter.Check(key, interval, max)
	assert.Nil(t, err)
	assert.Equal(t, ok, false)

	time.Sleep(10 * time.Millisecond)

	ok, err = limiter.Check(key, interval, max)
	assert.Nil(t, err)
	assert.Equal(t, ok, true)
}
