package middleware

import (
	"net/http"
	"time"

	"todo-api/pkg/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Info("Request:" + r.RemoteAddr + " " + r.Method + " " + r.URL.Path + " Duration: " + duration.String())
	})
}
