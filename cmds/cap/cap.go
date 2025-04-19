package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "cap")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, logger); err != nil {
		logger.Error("run", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	logger.Info("Hello!",
		"GOMAXPROCS", runtime.GOMAXPROCS(0),
		"GOARCH", runtime.GOARCH,
		"GOOS", runtime.GOOS,
		"GOLANG_VERSION", runtime.Version(),
	)
	defer logger.Info("Bye!")

	server := NewServer(ctx, logger)
	go shutdown(ctx, logger, server)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func shutdown(ctx context.Context, logger *slog.Logger, server *http.Server) {
	<-ctx.Done()

	logger.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		logger.Error("Failed to shutdown server", "error", err)
	}
	logger.Info("Server shutdown")
}
