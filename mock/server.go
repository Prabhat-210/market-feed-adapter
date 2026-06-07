package mock

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Server is a mock WebSocket server that mimics Upstox V3 feed
type Server struct {
	addr       string
	intervalMs int
	log        zerolog.Logger
}

func NewServer(addr string, intervalMs int, log zerolog.Logger) *Server {
	return &Server{
		addr:       addr,
		intervalMs: intervalMs,
		log:        log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// mock authorize endpoint — returns fake wss:// URL
	mux.HandleFunc("/v3/feed/market-data-feed/authorize", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "success",
			"data": {
				"authorized_redirect_uri": "ws://` + s.addr + `/feed"
			}
		}`))
	})

	// mock WebSocket feed endpoint
	mux.HandleFunc("/feed", s.handleFeed)

	srv := &http.Server{Addr: s.addr, Handler: mux}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	s.log.Info().Str("addr", s.addr).Msg("mock feed server started")
	return srv.ListenAndServe()
}
func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to upgrade connection")
		return
	}
	defer conn.Close()

	s.log.Info().Str("remote", r.RemoteAddr).Msg("mock client connected")

	gen := NewGenerator()
	ticker := time.NewTicker(time.Duration(s.intervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			payload, err := gen.Next()
			if err != nil {
				s.log.Error().Err(err).Msg("failed to generate mock tick")
				continue
			}

			// print what we are sending
			var pretty map[string]interface{}
			if err := json.Unmarshal(payload, &pretty); err == nil {
				s.log.Debug().
					RawJSON("payload", payload).
					Msg("sending mock tick")
			}

			if err := conn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
				s.log.Warn().Err(err).Msg("client disconnected")
				return
			}
		}
	}
}
