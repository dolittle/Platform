package cmd

import (
	"time"

	"github.com/dolittle/platform-router/admin"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/http"
	"github.com/dolittle/platform-router/kubernetes"
	"github.com/dolittle/platform-router/microservices"
	"github.com/dolittle/platform-router/proxy"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Starts the proxy server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		config, logger, err := config.SetupFor(cmd)
		if err != nil {
			return err
		}
		log.Logger = logger

		registry := microservices.NewRegistry()

		client, err := kubernetes.NewClient()
		if err != nil {
			return err
		}

		converter := &microservices.Converter{
			Config: config,
		}
		go converter.WatchConfig(cmd.Context())

		podWatcher := &kubernetes.PodWatcher{
			Client:                  client,
			Config:                  config,
			LabelSelectorConfigPath: "kubernetes.label-selector",
			Handler: &microservices.Updater{
				Registry:  registry,
				Converter: converter,
			},
		}
		go podWatcher.Run(time.Minute, cmd.Context())

		router := mux.NewRouter()
		admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry, config)
		proxy.AddApi(router.PathPrefix("/proxy").Subrouter(), registry, config, cmd.Context())

		server := &http.ReloadingServer{
			Handler:        router,
			WriteTimeout:   15 * time.Second,
			ReadTimeout:    15 * time.Second,
			Config:         config,
			PortConfigPath: "proxy.port",
		}

		go server.ListenAndServe()
		<-cmd.Context().Done()
		server.Shutdown(10 * time.Second)
		return nil
	},
}

func init() {
	proxyCmd.Flags().Int("proxy.port", 8080, "The port the proxy server should listen on")
	proxyCmd.Flags().String("proxy.tenant-header", "Tenant-ID", "The name of the header to use to resolve the request Tenant-ID")
	proxyCmd.Flags().String("kubernetes.label-selector", "tenant,application,environment,microservice,!infrastructure", "The label selector that will be used by the Kubernetes informer to only watch relevant Microservice pods")
}
