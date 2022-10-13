package proxy

import "github.com/dolittle/platform-router/microservices"

type Route struct {
	Host string
	Path string
}

type Paths map[string][]microservices.Port
