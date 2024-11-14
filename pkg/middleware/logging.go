package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ip := r.RemoteAddr
			userAgent := r.UserAgent()

			next.ServeHTTP(w, r)

			logger.Printf("Completed %s %s from %s [%s] in %v", r.Method, r.URL.Path, ip, userAgent, time.Since(start))
		})
	}
}
