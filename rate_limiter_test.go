package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimiterFull(t *testing.T) {

	ratelimiter := NewRateLimiter(3, time.Duration(time.Second*10))

	var wg sync.WaitGroup

	past := time.Now()
	for i := 0; i < 400; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ratelimiter.Push(time.Now())

		}(i)
	}
	wg.Wait()

	if !ratelimiter.LimitReached() {
		t.Error("rate limiter not full")
	}

	if time.Since(past) > time.Duration(time.Second*10) {
		t.Error("total time exceeded limiter duration")
	}

}
