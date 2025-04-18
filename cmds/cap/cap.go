package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "cap")
	ctx := context.Background()

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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	return nil
}
