package main

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start timer
		startTime := time.Now()

		// Process request
		next.ServeHTTP(w, r)

		// Calculate duration
		duration := time.Now().Sub(startTime)

		// Log duration
		log.Printf("[%s] %s %s %v", r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}
