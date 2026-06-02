package main

import (
	"context"
	"feed-adapter/internal/bootstrap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := bootstrap.Initialize()
	if err != nil {
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		os.Exit(1)
	}
	// app.Close()//close pipeline
}
