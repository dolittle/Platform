package proxy

import (
	"fmt"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"strings"
)

type RouterHandler struct {
	Registry *microservices.Registry
	Resolver *PortResolver
	Proxy    *httputil.ReverseProxy
	Config   *config.Config
}

func (rh RouterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenant := r.Header.Get(rh.Config.String("proxy.tenant-header"))
	application, environment, microservice, portName := getPathVars(r)

	logger := log.With().
		Str("component", "RouterHandler").
		Str("method", r.Method).
		Str("host", r.URL.Host).
		Str("path", r.URL.Path).
		Logger()

	logger.Debug().Msg("Handling request")

	pathPrefix, err := getActualPath(r)
	if err != nil {
		http.Error(w, "Could not parse path", http.StatusInternalServerError)
		logger.Debug().Err(err).Msg("Failed to get actual path")
		return
	}

	pathSuffix := strings.TrimPrefix(r.URL.Path, pathPrefix)
	if pathSuffix == r.URL.Path {
		http.Error(w, "Could not parse path", http.StatusInternalServerError)
		logger.Debug().Msg("Failed to remove path prefix")
		return
	}

	id := microservices.ToIdentity(tenant, application, environment, microservice)
	ms, found := rh.Registry.Get(id)
	if !found {
		http.NotFound(w, r)
		logger.Debug().Interface("id", id).Msg("Could not find microservice")
		return
	}

	port, found := rh.Resolver.ResolvePort(portName, ms.Ports)
	if !found {
		http.NotFound(w, r)
		logger.Debug().Str("port", portName).Msg("Could not find port")
		return
	}

	host := fmt.Sprintf("%s:%d", ms.IP, port)
	path := fmt.Sprintf("/%s", pathSuffix)

	route := Route{
		Host: host,
		Path: path,
	}

	rh.Proxy.ServeHTTP(w, RequestWithRoute(r, route))
}

const RouterPath = "/{application}/{environment}/{microservice}/{port}/"

func getPathVars(r *http.Request) (application, environment, microservice, port string) {
	vars := mux.Vars(r)
	return vars["application"], vars["environment"], vars["microservice"], vars["port"]
}

func getActualPath(r *http.Request) (string, error) {
	application, environment, microservice, port := getPathVars(r)
	url, err := mux.CurrentRoute(r).URLPath("application", application, "environment", environment, "microservice", microservice, "port", port)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
