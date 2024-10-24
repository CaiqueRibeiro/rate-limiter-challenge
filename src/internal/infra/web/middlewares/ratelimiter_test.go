package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RateLimiterMock struct {
	mock.Mock
}

func (m *RateLimiterMock) Check(ctx context.Context, r *http.Request) (*strategies.LimitResponse, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(*strategies.LimitResponse), args.Error(1)
}

func TestRateLimiterMiddlewareHandleAllow(t *testing.T) {
	mockLimiter := new(RateLimiterMock)
	middleware := NewRateLimiterMiddleware(mockLimiter)

	mockLimiter.On("Check", mock.Anything, mock.Anything).Return(&strategies.LimitResponse{
		Result:    strategies.Allow,
		Limit:     10,
		Remaining: 5,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	req := httptest.NewRequest(http.MethodGet, "/some-endpoint", nil)
	rr := httptest.NewRecorder()

	middleware.Handle(nextHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "10", rr.Header().Get("X-RateLimit-Limit"))
	assert.Equal(t, "5", rr.Header().Get("X-RateLimit-Remaining"))
	assert.NotEmpty(t, rr.Header().Get("X-RateLimit-Reset"))
	assert.Equal(t, "success", rr.Body.String())

	mockLimiter.AssertExpectations(t)
}

func TestRateLimiterMiddlewareHandleDeny(t *testing.T) {
	mockLimiter := new(RateLimiterMock)
	middleware := NewRateLimiterMiddleware(mockLimiter)

	mockLimiter.On("Check", mock.Anything, mock.Anything).Return(&strategies.LimitResponse{
		Result:    strategies.Deny,
		Limit:     int64(10),
		Remaining: int64(0),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	middleware.Handle(http.NotFoundHandler()).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusTooManyRequests, rr.Code)
	assert.Contains(t, rr.Body.String(), "maximum number of requests")
	mockLimiter.AssertExpectations(t)
}

func TestRateLimiterMiddlewareHandleInternalServerError(t *testing.T) {
	mockLimiter := new(RateLimiterMock)
	middleware := NewRateLimiterMiddleware(mockLimiter)

	mockLimiter.On("Check", mock.Anything, mock.Anything).Return((*strategies.LimitResponse)(nil), http.ErrHandlerTimeout)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	middleware.Handle(http.NotFoundHandler()).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), http.ErrHandlerTimeout.Error())
	mockLimiter.AssertExpectations(t)
}
