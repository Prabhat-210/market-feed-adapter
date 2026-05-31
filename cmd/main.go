package main

import (
	"context"
	"feed-adapter/internal/bootstrap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_, err := bootstrap.Initialize()
	if err != nil {
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-ctx.Done()

	// app.Close()//close pipeline
}
