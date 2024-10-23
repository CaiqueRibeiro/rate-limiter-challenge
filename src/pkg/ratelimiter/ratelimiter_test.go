package ratelimiter

import (
	"context"
	"errors"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type StrategyMock struct {
	mock.Mock
}

func (m *StrategyMock) CheckLimit(ctx context.Context, r *strategies.Request) (*strategies.LimitResponse, error) {
	args := m.Called(ctx, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*strategies.LimitResponse), args.Error(1)
}

func (m *StrategyMock) CheckTokenLimit(ctx context.Context, token string) (int64, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}
	return args.Get(0).(int64), args.Error(1)
}

func TestRateLimiterByIP(t *testing.T) {
	strategyMock := new(StrategyMock)
	ipMaxReqs := 5
	timeWindow := 1000
	limiter := NewRateLimiter(strategyMock, ipMaxReqs, timeWindow)

	t.Run("Should deny the request", func(t *testing.T) {
		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		request := strategies.Request{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		response := strategies.LimitResponse{
			Result:    strategies.Deny,
			Limit:     int64(ipMaxReqs),
			Total:     int64(ipMaxReqs),
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(timeWindow) * time.Millisecond),
		}

		strategyMock.On("CheckLimit", ctx, &request).Return(&response, nil)

		result, err := limiter.Check(ctx, r)

		assert.Nil(t, err)
		assert.Equal(t, response, *result)
		strategyMock.AssertExpectations(t)

		strategyMock.ExpectedCalls = nil
	})

	t.Run("Should allow the request", func(t *testing.T) {
		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		request := strategies.Request{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		response := strategies.LimitResponse{
			Result:    strategies.Allow,
			Limit:     int64(ipMaxReqs),
			Total:     int64(ipMaxReqs - 1),
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(timeWindow) * time.Millisecond),
		}

		strategyMock.On("CheckLimit", ctx, &request).Return(&response, nil)

		result, err := limiter.Check(ctx, r)

		assert.Nil(t, err)
		assert.Equal(t, response, *result)
		strategyMock.AssertExpectations(t)

		strategyMock.ExpectedCalls = nil
	})

	t.Run("Should return error if unexpected error occurs in strategy", func(t *testing.T) {
		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		request := strategies.Request{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(ipMaxReqs),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		strategyMock.On("CheckLimit", ctx, &request).Return(nil, errors.New("error-by-redis-limiter"))

		result, err := limiter.Check(ctx, r)

		assert.Error(t, err, "error-by-redis-limiter")
		assert.Nil(t, result)
		strategyMock.AssertExpectations(t)
	})
}

func TestRateLimiterByToken(t *testing.T) {
	strategyMock := new(StrategyMock)
	ipMaxReqs := 5
	timeWindow := 1000
	limiter := NewRateLimiter(strategyMock, ipMaxReqs, timeWindow)

	t.Run("Should deny the request", func(t *testing.T) {
		token := "dummy_token"
		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("API_KEY", token)

		request := strategies.Request{
			Key:      token,
			Limit:    int64(50),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		response := strategies.LimitResponse{
			Result:    strategies.Deny,
			Limit:     int64(50),
			Total:     50,
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(timeWindow) * time.Millisecond),
		}

		strategyMock.On("CheckTokenLimit", ctx, token).Return(int64(50), nil)
		strategyMock.On("CheckLimit", ctx, &request).Return(&response, nil)

		result, err := limiter.Check(ctx, r)

		assert.Nil(t, err)
		assert.Equal(t, response, *result)
		strategyMock.AssertExpectations(t)

		strategyMock.ExpectedCalls = nil
	})

	t.Run("Should allow the request", func(t *testing.T) {
		token := "dummy_token"
		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("API_KEY", token)

		request := strategies.Request{
			Key:      token,
			Limit:    int64(50),
			Duration: time.Duration(timeWindow) * time.Millisecond,
		}

		response := strategies.LimitResponse{
			Result:    strategies.Allow,
			Limit:     int64(50),
			Total:     1,
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(timeWindow) * time.Millisecond),
		}

		strategyMock.On("CheckTokenLimit", ctx, token).Return(int64(50), nil)
		strategyMock.On("CheckLimit", ctx, &request).Return(&response, nil)

		result, err := limiter.Check(ctx, r)

		assert.Nil(t, err)
		assert.Equal(t, response, *result)
		strategyMock.AssertExpectations(t)

		strategyMock.ExpectedCalls = nil
	})
}
