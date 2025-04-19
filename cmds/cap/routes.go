package main

import (
	"context"
	"log/slog"
	"net/http"
)

func NewServer(ctx context.Context, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	return &http.Server{
		Addr:    ":8080",
		Handler: logMiddleware(recoverMiddleware(mux, logger), logger),
	}
}

func logMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func recoverMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recovered from panic", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
