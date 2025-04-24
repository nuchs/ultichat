package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer(ctx context.Context, logger *slog.Logger) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", serverWS(ctx, logger))

	return &http.Server{
		Addr:    ":8080",
		Handler: logMiddleware(recoverMiddleware(mux, logger), logger),
	}
}

func serverWS(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("Failed to upgrade connection", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for {
			if ctx.Err() != nil {
				return
			}
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				logger.Error("Failed to read message", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			logger.Info("Received message", "type", msgType, "message", string(msg))
			if err := ws.WriteMessage(websocket.TextMessage, []byte("boop")); err != nil {
				logger.Error("Failed to send message", "error", err)
			}
		}
	})
}

func logMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		logger.Info("Finished request", "method", r.Method, "path", r.URL.Path, "status", fmt.Sprintf("%+v", w.Header()))
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
