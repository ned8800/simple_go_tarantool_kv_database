package middleware

import (
	"context"
	"net/http"
)

// AccessLogMiddleware logs any request
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("Req-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("Req-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
