package proxy

import (
	"context"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"net/http/httputil"
)

func AddApi(router *mux.Router, registry *microservices.Registry, config *config.Config, ctx context.Context) {
	resolver := &PortResolver{
		Config: config,
	}
	router.PathPrefix(RouterPath).Handler(RouterHandler{
		Registry: registry,
		Resolver: resolver,
		Proxy: &httputil.ReverseProxy{
			Director: director,
		},
		Config: config,
	})
	go resolver.WatchConfig(ctx)
}
