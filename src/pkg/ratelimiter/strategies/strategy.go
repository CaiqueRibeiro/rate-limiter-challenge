package strategies

import (
	"context"
	"time"
)

type Result int

const (
	Allow Result = 1
	Deny  Result = -1
)

type Request struct {
	Key      string
	Limit    int64
	Duration time.Duration
}

type Response struct {
	Result    Result
	Limit     int64
	Total     int64
	Remaining int64
	ExpiresAt time.Time
}

type LimiterStrategyInterface interface {
	Check(ctx context.Context, r *Request) (*Response, error)
}
