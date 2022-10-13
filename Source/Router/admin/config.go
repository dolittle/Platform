package admin

import (
	"fmt"
	"github.com/dolittle/platform-router/config"
	"net/http"
)

type ConfigHandler struct {
	Config *config.Config
}

func (ch ConfigHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	data, err := ch.Config.MarshalYaml()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal config as YAML: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
