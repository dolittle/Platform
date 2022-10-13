package cmd

import (
	"fmt"
	"github.com/dolittle/platform-router/admin"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"time"
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
		//ctx, stop := context.WithCancel(cmd.Context())
		//ctx := cmd.Context()

		registry := microservices.NewRegistry()
		// TODO: Remove after testing
		registry.Upsert(microservices.Microservice{
			Identity: microservices.ToIdentity("", "a", "b", "c"),
			IP:       "127.0.0.1",
			Ports: map[microservices.Port]int{
				microservices.Port{
					Container: "runtime",
					Port:      "http",
				}: 6006,
			},
		})

		router := mux.NewRouter()
		admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry, config)
		//proxy.AddApi(router.PathPrefix("/proxy").Subrouter(), registry, config)

		server := &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%d", config.Int("proxy.port")),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		log.Info().Str("address", server.Addr).Msg("Starting server")
		return server.ListenAndServe()

		//log.Info().Msg(config.String("listen-on"))
		//return nil

		//log.Info().Str("address", config.String("listenOn")).Msg("Starting server")
		//// microservicesConfig, err := microservices.GetConfiguration(viper.GetViper())
		//microservicesConfig := microservices.Configuration{}
		//err = config.Unmarshal("microservices", &microservicesConfig)
		//if err != nil {
		//	return err
		//}

		//log.Info().Interface("microservicesConfig", microservicesConfig).Msg("TEST")

		//client, err := kubernetes.NewClient()
		//if err != nil {
		//	return err
		//}

		//kubernetes.StartNewPodWatcher(
		//	client,
		//	time.Minute,
		//	"tenant,application,environment,microservice,!infrastructure",
		//	microservices.NewUpdater(registry, microservicesConfig.Microservice),
		//	cmd.Context().Done(),
		//)

		//router := mux.NewRouter()
		//admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry)
		//microservices.AddProxy(router, registry, microservicesConfig.Proxy)
		//server := &http.Server{
		//	Handler:      router,
		//	Addr:         viper.GetString("listenOn"),
		//	WriteTimeout: 15 * time.Second,
		//	ReadTimeout:  15 * time.Second,
		//}

		//go server.ListenAndServe()
		//<-cmd.Context().Done()
		//server.Shutdown(cmd.Context())
		//return nil
	},
}

func init() {
	proxyCmd.Flags().Int("proxy.port", 8080, "The port the proxy server should listen on")
}
