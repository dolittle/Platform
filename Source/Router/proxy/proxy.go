package proxy

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

func director(r *http.Request) {
	route := GetRequestRoute(r)

	log.Trace().
		Str("host", route.Host).
		Str("path", route.Path).
		Msg("Routing request to")

	r.URL.Scheme = "http"
	r.URL.Host = route.Host
	r.URL.Path = route.Path
}
