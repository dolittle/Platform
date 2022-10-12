package proxy

import (
	"fmt"
	"github.com/dolittle/platform-router/microservices"
	"github.com/knadh/koanf"
)

type PortResolver struct {
	Config *koanf.Koanf
}

func (pr *PortResolver) ResolvePort(port string, ports microservices.Ports) (int, bool) {
	path := fmt.Sprintf("paths.%s", port)
	if !pr.Config.Exists(path) {
		return 0, false
	}

	var possiblePorts []microservices.Port
	if err := pr.Config.Unmarshal(path, &possiblePorts); err != nil {
		return 0, false
	}

	for _, possiblePort := range possiblePorts {
		port, found := ports[possiblePort]
		if found {
			return port, true
		}
	}

	return 0, false
}
