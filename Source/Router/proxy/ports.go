package proxy

import (
	"context"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/microservices"
	"github.com/rs/zerolog/log"
)

type PortResolver struct {
	Config          *config.Config
	configuredPaths Paths
}

func (pr *PortResolver) ResolvePort(port string, ports microservices.Ports) (int, bool) {
	possiblePorts := pr.configuredPaths[port]

	for _, possiblePort := range possiblePorts {
		port, found := ports[possiblePort]
		if found {
			return port, true
		}
	}

	return 0, false
}

func (pr *PortResolver) WatchConfig(ctx context.Context) {
	for {
		changed := pr.Config.Changed()

		newPaths := Paths{}
		if err := pr.Config.Unmarshal("proxy.paths", &newPaths); err != nil {
			log.Warn().Err(err).Msg("Failed to unmarshal proxy paths")
		} else {
			pr.configuredPaths = newPaths
		}

		select {
		case <-changed:
			continue
		case <-ctx.Done():
			break
		}
	}
}
