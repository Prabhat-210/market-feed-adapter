package pipeline

import (
	"context"

	"feed-adapter/internal/model"
	"feed-adapter/internal/platform/config"

	"github.com/rs/zerolog"
)

type Pipeline struct {
	decoderChannels     []chan *model.RawTick
	interpreterChannels []chan *model.DecodedTick
	log                 zerolog.Logger
	cfg                 *config.Config
}

func NewPipeline(cfg *config.Config, log zerolog.Logger) *Pipeline {

	decoderChans := make([]chan *model.RawTick, cfg.DecoderCount)
	for i := range decoderChans {
		decoderChans[i] = make(chan *model.RawTick, cfg.BufferSize)
	}

	interpreterChans := make([]chan *model.DecodedTick, cfg.InterpreterCount)
	for i := range interpreterChans {
		interpreterChans[i] = make(chan *model.DecodedTick, cfg.BufferSize)
	}
	return &Pipeline{
		decoderChannels:     decoderChans,
		interpreterChannels: interpreterChans,
		log:                 log,
		cfg:                 cfg,
	}
}

func (p *Pipeline) Start(ctx context.Context, in <-chan *model.RawTick) {
	ingestor := NewIngestor(in, p.decoderChannels, p.log)
	go ingestor.Start(ctx)
}
