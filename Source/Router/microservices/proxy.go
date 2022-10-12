package microservices

import (
	"bytes"
	// "context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var (
	ErrRequestMissingTenant = func(tenantIDHeader string) error {
		return fmt.Errorf("request to proxy is missing tenant-id on header %s", tenantIDHeader)
	}
	ErrRequestHasMultipleTenants = func(tenantIDHeader string) error {
		return fmt.Errorf("request to proxy has multiple tenant-ids on header %s", tenantIDHeader)
	}
	ErrRequestMissingApplication  = errors.New("request to proxy is missing application-id variable")
	ErrRequestMissingEnvironment  = errors.New("request to proxy is missing environment variable")
	ErrRequestMissingMicroservice = errors.New("request to proxy is missing microservice-id variable")
)

var (
	applicationIDVariable  = "applicationId"
	environmentVariable    = "environment"
	microserviceIDVariable = "microserviceId"
	proxyRoute             = fmt.Sprintf("/{%s}/{%s}/{%s}", applicationIDVariable, environmentVariable, microserviceIDVariable)
)
var internalErrorHeader = "X-Routing-Error"

type roundTripperFilter struct {
	parent http.RoundTripper
}

func (rtf *roundTripperFilter) RoundTrip(r *http.Request) (*http.Response, error) {
	if err, ok := r.Header[internalErrorHeader]; ok {
		response := http.Response{
			StatusCode: http.StatusInternalServerError,
			Request:    r,
			Body:       ioutil.NopCloser(bytes.NewBufferString(strings.Join(err, ",")))}
		return &response, nil
	}
	return rtf.parent.RoundTrip(r)
}

func AddProxy(router *mux.Router, registry *Registry, config ProxyConfiguration) {
	router.Handle(proxyRoute, newProxy(registry, config))
}

func newProxy(registry *Registry, config ProxyConfiguration) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: director(registry, config),
		Transport: &roundTripperFilter{
			parent: http.DefaultTransport,
		},
	}
}

// func middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Parse input

// 		ctx := r.Context()
// 		ctx = context.WithValue(ctx, "host", "")
// 		ctx = context.WithValue(ctx, "port", "")
// 		ctx = context.WithValue(ctx, "proxy-path", "")

// 		if false {
// 			w.WriteHeader(500)
// 			return
// 		}

// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

func director(registry *Registry, config ProxyConfiguration) func(r *http.Request) {
	return func(request *http.Request) {
		// tenantID := request.Context().Value("tenant-id")
		logger := log.With().Str("requestURL", request.URL.String()).Logger()
		request.URL.Scheme = "http"
		logger.Info().Str("host", request.URL.Host).Msg("The host")
		microserviceIdentity, err := parseMicroserviceIdentity(request, config)
		if err != nil {
			logger.Error().Err(err).Interface("microserviceIdentity", microserviceIdentity).Msg("Failed to parse request to proxy")
			request.Header.Set(internalErrorHeader, fmt.Sprintf("Failed to parse request to proxy. %s", err.Error()))
			return
		}

		microservice, ok := registry.Get(microserviceIdentity)
		if !ok {
			logger.Error().Err(err).Interface("microserviceIdentity", microserviceIdentity).Msg("No microservice with identity in registry")
			request.Header.Set(internalErrorHeader, fmt.Sprintf("No microservice with identity %+v in registry", microserviceIdentity))
			return
		}
		logger = logger.With().Interface("microservice", microservice).Logger()
		proxyPath := strings.TrimPrefix(request.URL.Path, fmt.Sprintf("/%s/%s/%s", microservice.Identity.Application, microservice.Identity.Environment, microservice.Identity.Microservice))
		logger.Info().Str("path", request.URL.Path).Str("proxyPath", proxyPath).Msg("TEST")
		// host := fmt.Sprintf("%s:%d", microservice.IP, config.TerminalPort)
		// logger.Debug().Str("newHost", host).Msg("Proxying request")
		// request.URL.Host = host
	}
}

func parseMicroserviceIdentity(request *http.Request, config ProxyConfiguration) (Identity, error) {
	identity := Identity{}
	vars := mux.Vars(request)

	tenantIDs, ok := request.Header[config.TenantIDHeader]
	if !ok {
		return identity, ErrRequestMissingTenant(config.TenantIDHeader)
	}

	if len(tenantIDs) > 1 {
		return identity, ErrRequestHasMultipleTenants(config.TenantIDHeader)
	}

	identity.Tenant = identityFromString[TenantID](tenantIDs[0])
	applicationID, ok := vars[applicationIDVariable]
	if !ok {
		return identity, ErrRequestMissingApplication
	}

	identity.Application = identityFromString[ApplicationID](applicationID)
	environment, ok := vars[environmentVariable]
	if !ok {
		return identity, ErrRequestMissingEnvironment
	}

	identity.Environment = identityFromString[Environment](environment)
	microserviceID, ok := vars[microserviceIDVariable]
	if !ok {
		return identity, ErrRequestMissingMicroservice
	}

	identity.Microservice = identityFromString[MicroserviceID](microserviceID)
	return identity, nil
}
