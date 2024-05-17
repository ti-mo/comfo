package main

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger is an http.Handler wrapper that generates HTTP logs on stdout.
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		slog.Debug("Request",
			slog.String("method", r.Method),
			slog.String("URI", r.RequestURI),
			slog.String("route", name),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// Headers sets response headers on responses served from the API.
func Headers(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		inner.ServeHTTP(w, r)
	})
}
