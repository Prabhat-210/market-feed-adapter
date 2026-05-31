package bootstrap

import (
	"context"
	"feed-adapter/internal/config"
	"feed-adapter/internal/platform/logger"

	"github.com/rs/zerolog"
)

type Application struct {
	log zerolog.Logger
	cfg *config.Config
}

func Initialize() (*Application, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	log := logger.NewLogger(cfg.ServiceName, cfg.Environment, cfg.LogLevel)

	return &Application{
		log: log,
		cfg: cfg,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	a.log.Info().Msg("application started")

	<-ctx.Done()

	a.log.Info().Msg("application stopped")

	return nil
}
