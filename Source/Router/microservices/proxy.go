package microservices

import (
	"bytes"
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
	ErrRequestMissingTenant       = errors.New("request to proxy is missing tenant-id variable")
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
			StatusCode: http.StatusBadRequest,
			Request:    r,
			Body:       ioutil.NopCloser(bytes.NewBufferString(strings.Join(err, ",")))}
		return &response, nil
	}
	return rtf.parent.RoundTrip(r)
}

func AddProxy(router *mux.Router, registry *Registry) {
	router.Handle(proxyRoute, newProxy(registry))
}

func newProxy(registry *Registry) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: director(registry),
		Transport: &roundTripperFilter{
			parent: http.DefaultTransport,
		},
	}
}

func director(registry *Registry) func(r *http.Request) {
	return func(r *http.Request) {
		logger := log.With().Str("requestURL", r.URL.String()).Logger()
		r.URL.Scheme = "http"
		logger.Info().Str("host", r.URL.Host).Msg("The host")
		microserviceIdentity, err := parseMicroserviceIdentity(r)
		if err != nil {
			logger.Error().Err(err).Interface("microserviceIdentity", microserviceIdentity).Msg("Failed to parse request to proxy")
			r.Header.Set(internalErrorHeader, fmt.Sprintf("Failed to parse request to proxy. %s", err.Error()))
			return
		}

		microservice, ok := registry.Get(microserviceIdentity)
		if !ok {
			logger.Error().Err(err).Interface("microserviceIdentity", microserviceIdentity).Msg("No microservice with identity in registry")
			r.Header.Set(internalErrorHeader, fmt.Sprintf("No microservice with identity %+v in registry", microserviceIdentity))
			return
		}

		logger = logger.With().Interface("microservice", microservice).Logger()
		port := 8080 // TODO:: Figure out the port
		host := fmt.Sprintf("%s:%d", microservice.IP, port)
		logger.Debug().Str("newHost", host).Msg("Proxying request")
		r.URL.Host = host
	}
}

func parseMicroserviceIdentity(request *http.Request) (Identity, error) {
	identity := Identity{}
	vars := mux.Vars(request)

	identity.Tenant = identityFromString[TenantID]("some-tenant-id") // TODO: Where to get this from. [Tenant-Id]?
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
