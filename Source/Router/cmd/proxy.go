package cmd

import (
	"github.com/dolittle/platform-router/admin"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/kubernetes"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var proxy = &cobra.Command{
	Use:   "proxy",
	Short: "Starts the proxy server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		config, logger, err := config.SetupFor(cmd)
		if err != nil {
			return err
		}
		log.Logger = logger

		log.Info().Msg(config.String("listen-on"))
		return nil

		log.Info().Str("address", config.String("listenOn")).Msg("Starting server")
		// microservicesConfig, err := microservices.GetConfiguration(viper.GetViper())
		microservicesConfig := microservices.Configuration{}
		err = config.Unmarshal("microservices", &microservicesConfig)
		if err != nil {
			return err
		}

		log.Info().Interface("microservicesConfig", microservicesConfig).Msg("TEST")

		client, err := kubernetes.NewClient()
		if err != nil {
			return err
		}

		registry := microservices.NewRegistry()
		kubernetes.StartNewPodWatcher(
			client,
			time.Minute,
			"tenant,application,environment,microservice,!infrastructure",
			microservices.NewUpdater(registry, microservicesConfig.Microservice),
			cmd.Context().Done(),
		)

		router := mux.NewRouter()
		admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry)
		microservices.AddProxy(router, registry, microservicesConfig.Proxy)
		server := &http.Server{
			Handler:      router,
			Addr:         viper.GetString("listenOn"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		go server.ListenAndServe()
		<-cmd.Context().Done()
		server.Shutdown(cmd.Context())
		return nil
	},
}

func init() {
	proxy.Flags().String("listen-on", "localhost:8080", "The address to listen on")
}
