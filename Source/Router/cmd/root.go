package cmd

import (
	"net/http"
	"time"

	"github.com/dolittle/platform-router/admin"
	"github.com/dolittle/platform-router/config"
	"github.com/dolittle/platform-router/kubernetes"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "router",
	Short: "Router",
	Long:  `Entry router`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: cmd.OutOrStdout()}).Level(zerolog.TraceLevel)

		conf, err := config.LoadConfigFor(cmd)
		if err != nil {
			return err
		}

		log.Info().Str("address", conf.String("listenOn")).Msg("Starting server")
		// microservicesConfig, err := microservices.GetConfiguration(viper.GetViper())
		microservicesConfig := microservices.Configuration{}
		err = conf.Unmarshal("microservices", &microservicesConfig)
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

func init() {
	rootCmd.PersistentFlags().StringSlice("config", nil, "A configuration file to load, can be specified multiple times")
	rootCmd.Flags().String("listenOn", "localhost:8080", "The address to listen on")

	// viper.AutomaticEnv()
	// viper.SetConfigName("config")
	// viper.SetConfigType("yaml")
	// viper.AddConfigPath(".")
	// viper.SetDefault("listenOn", "localhost:8080")

	// if err := viper.ReadInConfig(); err != nil {
	// 	panic(err)
	// }
}
