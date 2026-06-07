package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"feed-adapter/internal/platform/config"
	"feed-adapter/mock"

	"github.com/rs/zerolog"
)

func main() {
	// force mock mode — ignores whatever is in .env
	os.Setenv("MOCK_MODE", "true")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg, err := config.Load()
	if err != nil {
		os.Exit(1)
	}
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Str("service", "mock-feed-server").
		Logger()

	// only start mock server — nothing else
	mockSrv := mock.NewServer(cfg.MockAddr, cfg.MockIntervalMs, log)
	if err := mockSrv.Start(ctx); err != nil {
		log.Fatal().Err(err).Msg("mock server failed")
	}

	<-ctx.Done()
}
