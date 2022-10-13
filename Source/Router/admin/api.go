package admin

import (
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
)

func AddApi(router *mux.Router, registry *microservices.Registry, config *config.Config) {
	router.Path("/registry").Handler(RegistryHandler{
		Registry: registry,
	})
	router.Path("/config").Handler(ConfigHandler{
		Config: config,
	})
}
