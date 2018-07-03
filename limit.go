package main

import (
	"time"

	"golang.org/x/time/rate"
)

const (
	disableRateLimiter = true
)

var (
	Limit   = rate.Limit(time.Second / 120)
	Request = 2
)

var lim = rate.NewLimiter(Limit, Request)

// There was once a time where the repainting had to be rate
// limited. This currently happens when processing a request
// to repaint a section of the grid

func throttled() bool {
	if disableRateLimiter {
		return false
	}
	return !lim.Allow()
}
