package proxy

import (
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/knadh/koanf"
	"net/http/httputil"
)

func AddApi(router *mux.Router, registry *microservices.Registry, config *koanf.Koanf) {
	router.PathPrefix(RouterPath).Handler(RouterHandler{
		Registry: registry,
		Resolver: &PortResolver{
			Config: config,
		},
		Proxy: &httputil.ReverseProxy{
			Director: director,
		},
	})
}
