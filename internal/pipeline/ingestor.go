package pipeline

import (
	"context"
	"feed-adapter/internal/model"

	"github.com/rs/zerolog"
)

type Ingestor struct {
	in       <-chan *model.RawTick
	decoders []chan *model.RawTick
	current  int
	log      zerolog.Logger
}

func NewIngestor(in <-chan *model.RawTick, decoders []chan *model.RawTick, log zerolog.Logger) *Ingestor {
	return &Ingestor{
		in:       in,
		decoders: decoders,
		log:      log,
	}
}

func (i *Ingestor) Start(ctx context.Context) {
	i.log.Info().Int("Collector count", len(i.decoders)).Msgf("starting ingestor")

	for {
		select {
		case <-ctx.Done():
			i.log.Warn().Msg("Ingestor stopped")
		case tick, ok := <-i.in:
			if !ok {
				i.log.Warn().Msg("feed channel closed, ingestor stopped")
				return
			}
			i.dispatcher(tick)
		}
	}
}

func (i *Ingestor) dispatcher(tick *model.RawTick) {
	target := i.current % len(i.decoders)
	i.current++
	select {
	//send data in chan
	case i.decoders[target] <- tick:
		i.log.Debug().Int("Current channel", i.current).Msg("decoder updates")
	default:
		i.log.Warn().Msgf("Skipping updates as decoder channel %d is full", i.current)
	}
}
