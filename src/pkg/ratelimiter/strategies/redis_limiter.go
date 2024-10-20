package strategies

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	KeyWithoutTTL = -1
	KeyNotFound   = -2
)

type RedisLimiter struct {
	Client *redis.Client
	Now    func() time.Time
}

func NewRedisLimiter(
	client *redis.Client,
	now func() time.Time,
) *RedisLimiter {
	return &RedisLimiter{
		Client: client,
		Now:    now,
	}
}

func (rls *RedisLimiter) Check(ctx context.Context, r *Request) (*Response, error) {
	p := rls.Client.Pipeline()
	getResult := p.Get(ctx, r.Key)
	ttlResult := p.TTL(ctx, r.Key)

	if _, err := p.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var ttlDuration time.Duration

	ttl, err := ttlResult.Result()
	if err != nil || ttl == KeyWithoutTTL || ttl == KeyNotFound {
		ttlDuration = r.Duration

		if err := rls.Client.Expire(ctx, r.Key, r.Duration).Err(); err != nil {
			return nil, err
		}
	} else {
		ttlDuration = ttl
	}

	currentCount, err := getResult.Int64()
	if err != nil && errors.Is(err, redis.Nil) {
		currentCount = 0
	}

	expiresAt := rls.Now().Add(ttlDuration)

	if currentCount >= r.Limit {
		return &Response{
			Result:    Deny,
			Total:     currentCount,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	incrResult := rls.Client.Incr(ctx, r.Key)
	nextTotal, err := incrResult.Result()
	if err != nil {
		return nil, err
	}

	if nextTotal > r.Limit {
		return &Response{
			Result:    Deny,
			Total:     nextTotal,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	return &Response{
		Result:    Allow,
		Total:     nextTotal,
		Limit:     r.Limit,
		Remaining: r.Limit - nextTotal,
		ExpiresAt: expiresAt,
	}, nil
}
