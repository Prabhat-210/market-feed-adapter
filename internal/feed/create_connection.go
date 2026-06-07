package feed

import (
	"feed-adapter/internal/model"
	"feed-adapter/internal/platform/config"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Connection struct {
	authorizeURL   string
	accessToken    string
	instruments    []string
	mode           string
	out            chan *model.RawTick
	log            zerolog.Logger
	reconnectMaxMS int
	httpClient     *http.Client
}

func NewConnection(cfg config.FeedConfig,
	log zerolog.Logger,
) *Connection {
	return &Connection{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		authorizeURL:   cfg.AuthorizeURL,
		accessToken:    cfg.AccessToken,
		instruments:    cfg.Instruments,
		mode:           cfg.Mode,
		out:            make(chan *model.RawTick, cfg.ConnectionChanSize),
		log:            log,
		reconnectMaxMS: cfg.ReconnectMaxMs,
	}
}

func (c *Connection) Channel() <-chan *model.RawTick {
	return c.out
}
