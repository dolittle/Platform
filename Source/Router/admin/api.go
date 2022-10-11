package admin

import (
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
)

func AddApi(router *mux.Router, registry *microservices.Registry) {
	router.Path("/registry").Handler(RegistryHandler{
		Registry: registry,
	})
}
