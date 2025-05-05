package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RPSLimiter struct {
	mu               sync.Mutex
	tokens           int
	cap              int
	lastModifiedTime time.Time
}

func NewRPSLimiter(tokens int) *RPSLimiter {
	if tokens <= 0 {
		return nil
	}

	return &RPSLimiter{
		tokens: tokens,
		cap:    tokens,
	}
}

func (r *RPSLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	elapsed := now.Sub(r.lastModifiedTime)
	tokensToAdd := int(elapsed / (time.Second / time.Duration(r.cap)))

	if tokensToAdd > 0 {
		r.tokens = tokensToAdd + r.tokens
		if r.tokens > r.cap {
			r.tokens = r.cap
		}
		r.lastModifiedTime = now
	}

	fmt.Printf("tokens : %d\n", r.tokens)

	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

type RateLimiterMiddleware struct {
	limiter RPSLimiter
	handler http.Handler
}

func NewRateLimiterMiddleware(h http.Handler, rps int) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		limiter: *NewRPSLimiter(rps),
		handler: h,
	}
}

func (l *RateLimiterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !l.limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("too many requests"))
		return
	}
	l.handler.ServeHTTP(w, r)
}
