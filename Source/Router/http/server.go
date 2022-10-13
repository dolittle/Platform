package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dolittle/platform-router/config"
	"github.com/rs/zerolog/log"
)

// ReloadingServer represents a http server that reloads itself whenever Config changes.
type ReloadingServer struct {
	Handler        http.Handler
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	Config         *config.Config
	PortConfigPath string
	port           int
	server         *http.Server
	isShutdown     bool
}

// ListenAndServe starts a http.Server with the port configured in config.Config and listens to changes in the configuration.
// If the configured port changes the server will restart with the new configured port.
func (rs *ReloadingServer) ListenAndServe() {
	for {
		if rs.isShutdown {
			return
		}

		rs.port = rs.Config.Int(rs.PortConfigPath)
		rs.server = &http.Server{
			Handler:      rs.Handler,
			WriteTimeout: rs.WriteTimeout,
			ReadTimeout:  rs.ReadTimeout,
			Addr:         fmt.Sprintf(":%d", rs.port),
		}

		go rs.shutdownServerIfPortChanges()
		log.Info().
			Str("address", rs.server.Addr).
			Msg("Starting HTTP server")
		err := rs.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().
				Err(err).
				Str("address", rs.server.Addr).
				Msg("Failed to start HTTP server")
		}
	}
}

func (rs *ReloadingServer) Shutdown(timeout time.Duration) {
	log.Info().
		Str("address", rs.server.Addr).
		Msg("Shutting down HTTP server")
	rs.isShutdown = true
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = rs.server.Shutdown(ctx)
}

func (rs *ReloadingServer) shutdownServerIfPortChanges() {
	for {
		<-rs.Config.Changed()
		if rs.Config.Int(rs.PortConfigPath) != rs.port {
			break
		}
	}
	rs.server.Shutdown(context.Background())
}
