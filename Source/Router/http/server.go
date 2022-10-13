package http

import (
	"context"
	"fmt"
	"github.com/dolittle/platform-router/config"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type ReloadingServer struct {
	Handler        http.Handler
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	Config         *config.Config
	PortConfigPath string
	port           int
	server         *http.Server
	shutdown       bool
}

func (rs *ReloadingServer) ListenAndServe() {
	for {
		if rs.shutdown {
			return
		}

		change := rs.Config.Changed()
		rs.port = rs.Config.Int(rs.PortConfigPath)
		rs.server = &http.Server{
			Handler:      rs.Handler,
			WriteTimeout: rs.WriteTimeout,
			ReadTimeout:  rs.ReadTimeout,
			Addr:         fmt.Sprintf(":%d", rs.port),
		}

		go func() {
			for {
				<-change
				change = rs.Config.Changed()

				if rs.Config.Int(rs.PortConfigPath) != rs.port {
					break
				}
			}
			rs.server.Shutdown(context.Background())
		}()

		log.Info().Str("address", rs.server.Addr).Msg("Starting HTTP server")
		err := rs.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Str("address", rs.server.Addr).Msg("Failed to start HTTP server")
		}
	}
}

func (rs *ReloadingServer) Shutdown(timeout time.Duration) {
	log.Info().Str("address", rs.server.Addr).Msg("Shutting down HTTP server")
	rs.shutdown = true
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	_ = rs.server.Shutdown(ctx)
}
