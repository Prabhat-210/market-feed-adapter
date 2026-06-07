package feed

import (
	"context"
	"encoding/json"
	"feed-adapter/internal/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type subscribeMsg struct {
	Guid   string        `json:"guid"`
	Method string        `json:"method"`
	Data   subscribeData `json:"data"`
}

type subscribeData struct {
	Mode           string   `json:"mode"`
	InstrumentKeys []string `json:"instrumentKeys"`
}

func (c *Connection) Start(ctx context.Context) {
	backoff := 100 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			c.log.Info().Msg("feed connection stopped")
			close(c.out)
			return
		default:
			if err := c.connect(ctx); err != nil {
				c.log.Err(err).Msg("feed connection failed")

				select {
				//this case is for if app stops in middle of backoff execution,
				// W/O this we need to wait till backoff time
				case <-ctx.Done():
					close(c.out)
					return
				case <-time.After(backoff):
				}
				backoff *= 2

				//set backoff cap to max if its going more than max
				if backoff > time.Duration(c.reconnectMaxMS)*time.Millisecond {
					backoff = time.Duration(c.reconnectMaxMS) * time.Millisecond
				}
				continue
			}
			//reset backoff if success
			backoff = 100 * time.Millisecond
		}
	}
}

func (c *Connection) connect(ctx context.Context) error {
	wsURL, err := c.fetchAuthorizedURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get authorized URL: %w", err)
	}

	c.log.Info().Str("url", wsURL).Msg("connecting to feed")

	//creating persistent connection
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}

	defer conn.Close()

	c.log.Info().Msg("feed connected")

	if err := c.subscribe(conn); err != nil {
		return err
	}

	c.log.Info().Strs("instruments", c.instruments).Str("mode", c.mode).Msg("subscription successful")
	return c.readLoop(ctx, conn)
}

// fetches websocket URL from upstocks
func (c *Connection) fetchAuthorizedURL(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.authorizeURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization",
		"Bearer "+c.accessToken,
	)
	req.Header.Set("Accept",
		"application/json",
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authorize API returned status %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			AuthorizedRedirectURI string `json:"authorized_redirect_uri"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.AuthorizedRedirectURI, nil
}

func (c *Connection) subscribe(conn *websocket.Conn) error {
	// Tell Upstox which instruments we want to receive which kind of ticks
	msg := subscribeMsg{
		Guid:   "market-feed-adapter",
		Method: "sub",
		Data: subscribeData{
			Mode:           c.mode,
			InstrumentKeys: c.instruments,
		},
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.BinaryMessage, payload)
}

func (c *Connection) readLoop(ctx context.Context, conn *websocket.Conn) error {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// Application shutting down.
			if ctx.Err() != nil {
				return nil
			}
			// Connection lost.
			return fmt.Errorf("read failed: %w", err)
		}
		// Stage 0 only receives bytes.
		// Decoding and validation happen later in pipeline stages.
		tick := &model.RawTick{
			Data:       data,
			ReceivedAt: time.Now(),
		}
		select {
		// here we Push raw message to pipeline.
		case c.out <- tick:
		default:
			// Drop incoming tick instead of blocking websocket reads.
			c.log.Warn().Msg("feed output channel full, dropping incoming tick")
		}
	}
}
