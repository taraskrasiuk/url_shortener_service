package middlewares

import (
	"net/http"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(50, 3) // 50 rps

func ReqRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
