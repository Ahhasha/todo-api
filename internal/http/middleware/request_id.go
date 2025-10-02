package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get(RequestIDHeader)
		if requestId == "" {
			requestId = uuid.New().String()
		}
		w.Header().Set(RequestIDHeader, requestId)
		next.ServeHTTP(w, r)

	})
}
