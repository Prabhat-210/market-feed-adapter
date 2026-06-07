package bootstrap

import (
	"context"

	"feed-adapter/internal/feed"
	"feed-adapter/internal/pipeline"
	"feed-adapter/internal/platform/config"
	"feed-adapter/internal/platform/logger"

	"github.com/rs/zerolog"
)

type Application struct {
	log      zerolog.Logger
	cfg      *config.Config
	feedConn *feed.Connection
	pipeline *pipeline.Pipeline
}

func Initialize() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(cfg.ServiceName, cfg.Environment, cfg.LogLevel)
	log.Info().Any("Config", cfg).Msg("config loaded...")

	feedConn := feed.NewConnection(cfg.FeedConfig, log)
	log.Info().Msg("feed connection initialized...")

	p := pipeline.NewPipeline(cfg, log)

	return &Application{
		log:      log,
		cfg:      cfg,
		feedConn: feedConn,
		pipeline: p,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {

	a.log.Info().Msg("application started")

	go a.feedConn.Start(ctx)

	a.pipeline.Start(ctx, a.feedConn.Channel())

	<-ctx.Done()

	a.log.Info().Msg("application stopped")

	return nil
}
