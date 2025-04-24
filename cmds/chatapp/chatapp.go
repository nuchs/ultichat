package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

const serverUrl = "ws://localhost:8080"

func receiveHandler(conn *websocket.Conn, logger *slog.Logger) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Failed to read message", "error", err)
			}
			return
		}
		logger.Info("Received message", "type", msgType, "message", string(msg))
	}
}

func shutdown(ctx context.Context, conn *websocket.Conn, logger *slog.Logger) {
	<-ctx.Done()
	logger.Info("Shutting down")
	if err := conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
		logger.Error("Failed to close websocket", "err", err)
	}
	if err := os.Stdin.Close(); err != nil {
		logger.Error("Failed to close stdin", "err", err)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "chatapp")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, serverUrl, nil)
	if err != nil {
		logger.Error("Failed to connect to server", "url", serverUrl)
	}
	defer conn.Close()
	go receiveHandler(conn, logger)

	scanner := bufio.NewScanner(os.Stdin)
	go shutdown(ctx, conn, logger)
	for scanner.Scan() {
		logger.Info("--- newline ---")
		if scanner.Err() != nil {
			logger.Info("Exiting processing loop")
			break
		}
		line := scanner.Text()
		trim := strings.TrimSpace(line)

		if err := conn.WriteMessage(websocket.TextMessage, []byte(trim)); err != nil {
			logger.Error("Failed to send message", "err", err, "msg", trim)
		}
	}
}
