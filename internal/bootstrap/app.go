package bootstrap

import (
	"context"
	"feed-adapter/internal/config"
	"feed-adapter/internal/feed"
	"feed-adapter/internal/platform/logger"

	"github.com/rs/zerolog"
)

type Application struct {
	log      zerolog.Logger
	cfg      *config.Config
	feedConn *feed.Connection
}

func Initialize() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(cfg.ServiceName, cfg.Environment, cfg.LogLevel)

	feedConn := feed.NewConnection(cfg.FeedConfig, log)

	return &Application{
		log:      log,
		cfg:      cfg,
		feedConn: feedConn,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {

	a.log.Info().Msg("application started")

	go a.feedConn.Start(ctx)

	<-ctx.Done()

	a.log.Info().Msg("application stopped")

	return nil
}
