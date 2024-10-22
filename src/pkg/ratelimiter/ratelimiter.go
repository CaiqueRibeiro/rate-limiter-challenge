package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
	rip "github.com/vikram1565/request-ip"
)

type RateLimiterInterface interface {
	Check(ctx context.Context, r *http.Request) (*strategies.LimitResponse, error)
}

type RateLimiter struct {
	Strategy         strategies.LimiterStrategyInterface
	MaxRequestsPerIP int
	TimeWindowMillis int
}

func NewRateLimiter(
	strategy strategies.LimiterStrategyInterface,
	ipMaxReqs int,
	timeWindow int,
) *RateLimiter {
	return &RateLimiter{
		Strategy:         strategy,
		MaxRequestsPerIP: ipMaxReqs,
		TimeWindowMillis: timeWindow,
	}
}

func (rl *RateLimiter) Check(ctx context.Context, r *http.Request) (*strategies.LimitResponse, error) {
	var key string
	var limit int64
	duration := time.Duration(rl.TimeWindowMillis) * time.Millisecond

	apiKey := r.Header.Get("API_KEY")

	if apiKey != "" {
		fmt.Printf("TOKEN A SER VERIFICADO: %s\n", apiKey)
		tokenMaxRequests, err := rl.Strategy.CheckTokenLimit(r.Context(), apiKey)

		if err != nil {
			fmt.Println("Não achou, vai ser por IP")
			key = rip.GetClientIP(r)
			limit = int64(rl.MaxRequestsPerIP)
		} else { // if no token found, set as IP even with API_KEY present
			key = apiKey
			limit = tokenMaxRequests
		}
	} else {
		key = rip.GetClientIP(r)
		limit = int64(rl.MaxRequestsPerIP)
	}

	req := &strategies.Request{
		Key:      key,
		Limit:    limit,
		Duration: duration,
	}

	result, err := rl.Strategy.CheckLimit(r.Context(), req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
