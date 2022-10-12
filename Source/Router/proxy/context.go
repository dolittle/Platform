package proxy

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrContextNotSet = errors.New("the route context was not set on the request")
)

const routeContextKey = "dolittle.io/platform-router/route"

func RequestWithRoute(r *http.Request, route Route) *http.Request {
	ctx := context.WithValue(r.Context(), routeContextKey, route)
	return r.WithContext(ctx)
}

func GetRequestRoute(r *http.Request) Route {
	route, ok := r.Context().Value(routeContextKey).(Route)
	if !ok {
		panic(ErrContextNotSet)
	}
	return route
}
