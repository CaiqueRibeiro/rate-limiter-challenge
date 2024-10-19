package ratelimiter

import (
	"context"
	"net/http"
	"time"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
	rip "github.com/vikram1565/request-ip"
)

type RateLimiterInterface interface {
	Check(ctx context.Context, r *http.Request) (*strategies.Response, error)
}

type RateLimiter struct {
	Strategy            strategies.LimiterStrategyInterface
	MaxRequestsPerIP    int
	MaxRequestsPerToken int
	TimeWindowMillis    int
}

func NewRateLimiter(
	strategy strategies.LimiterStrategyInterface,
	ipMaxReqs int,
	tokenMaxReqs int,
	timeWindow int,
) *RateLimiter {
	return &RateLimiter{
		Strategy:            strategy,
		MaxRequestsPerIP:    ipMaxReqs,
		MaxRequestsPerToken: tokenMaxReqs,
		TimeWindowMillis:    timeWindow,
	}
}

func (rl *RateLimiter) Check(ctx context.Context, r *http.Request) (*strategies.Response, error) {
	var key string
	var limit int64
	duration := time.Duration(rl.TimeWindowMillis) * time.Millisecond

	apiKey := r.Header.Get("API_KEY")

	if apiKey != "" {
		key = apiKey
		limit = int64(rl.MaxRequestsPerToken)
	} else {
		key = rip.GetClientIP(r)
		limit = int64(rl.MaxRequestsPerIP)
	}

	req := &strategies.Request{
		Key:      key,
		Limit:    limit,
		Duration: duration,
	}

	result, err := rl.Strategy.Check(r.Context(), req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
