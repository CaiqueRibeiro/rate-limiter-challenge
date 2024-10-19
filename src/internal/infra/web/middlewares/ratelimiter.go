package middlewares

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter"
	limiter "github.com/CaiqueRibeiro/rate-limiter-challenge/src/pkg/ratelimiter/strategies"
)

type RateLimiterMiddlewareInterface interface {
	Handle(next http.Handler) http.Handler
}

type RateLimiterMiddleware struct {
	Limiter ratelimiter.RateLimiterInterface
}

func NewRateLimiterMiddleware(
	limiter ratelimiter.RateLimiterInterface,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		Limiter: limiter,
	}
}

func (rlm *RateLimiterMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result, err := rlm.Limiter.Check(r.Context(), r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error(),
			})
			return
		}

		w.Header().Set("X-RateLimit-Limit", strconv.FormatInt(result.Limit, 10))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ExpiresAt.Unix(), 10))

		if result.Result == limiter.Deny {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "amount of requests exceeded",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
