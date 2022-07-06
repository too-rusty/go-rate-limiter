package ratelimiter

import (
	"time"
)

type RateLimiter struct {
	timestamps chan time.Time
	capacity   uint
	duration   time.Duration
}

// max request allowed capacity of the rate limiter and time period withing which all the requests hit for
func NewRateLimiter(capacity uint, duration time.Duration) *RateLimiter {
	ra := &RateLimiter{
		timestamps: make(chan time.Time, capacity),
		capacity:   capacity,
		duration:   duration,
	}
	go ra.consumer()
	return ra
}

// push all the timestamps here and let limiter do its job
func (ra *RateLimiter) Push(value time.Time) {
	go func() {
		ra.timestamps <- value
	}()
}

func (ra *RateLimiter) consumer() {
	for {
		tstamp := <-ra.timestamps
		now := time.Now()
		future := tstamp.Add(ra.duration)
		if future.After(now) {
			time.Sleep(future.Sub(now))
		}
	}

}

func (ra *RateLimiter) LimitReached() bool {
	return len(ra.timestamps) == int(ra.capacity)
	// https://groups.google.com/g/golang-nuts/c/L0wIBDr3HCc
}
