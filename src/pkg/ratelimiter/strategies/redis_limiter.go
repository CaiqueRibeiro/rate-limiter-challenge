package strategies

import (
	"context"
	"errors"
	"fmt"
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

func (rls *RedisLimiter) CheckTokenLimit(ctx context.Context, token string) (int64, error) {
	key := fmt.Sprintf("token_max_req:%s", token)
	getResult := rls.Client.Get(ctx, key)

	tokenMaxRequests, err := getResult.Int64()
	if err != nil && errors.Is(err, redis.Nil) {
		return 0, err
	}

	return tokenMaxRequests, nil
}

func (rls *RedisLimiter) CheckLimit(ctx context.Context, r *Request) (*LimitResponse, error) {
	key := fmt.Sprintf("limit:%s", r.Key)

	p := rls.Client.Pipeline()
	getResult := p.Get(ctx, key)
	ttlResult := p.TTL(ctx, key)

	if _, err := p.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var ttlDuration time.Duration

	ttl, err := ttlResult.Result()
	if err != nil || ttl == KeyWithoutTTL || ttl == KeyNotFound {
		ttlDuration = r.Duration

		if err := rls.Client.Expire(ctx, key, r.Duration).Err(); err != nil {
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
		return &LimitResponse{
			Result:    Deny,
			Total:     currentCount,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	incrResult := rls.Client.Incr(ctx, key)
	nextTotal, err := incrResult.Result()
	if err != nil {
		return nil, err
	}

	if nextTotal > r.Limit {
		return &LimitResponse{
			Result:    Deny,
			Total:     nextTotal,
			Limit:     r.Limit,
			Remaining: 0,
			ExpiresAt: expiresAt,
		}, nil
	}

	return &LimitResponse{
		Result:    Allow,
		Total:     nextTotal,
		Limit:     r.Limit,
		Remaining: r.Limit - nextTotal,
		ExpiresAt: expiresAt,
	}, nil
}
