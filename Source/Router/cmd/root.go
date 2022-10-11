package cmd

import (
	"net/http"
	"time"

	"github.com/dolittle/platform-router/admin"
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
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: cmd.OutOrStdout()}).Level(zerolog.DebugLevel)
		log.Info().Msg("Starting server")
		client, err := kubernetes.NewClient()
		if err != nil {
			return err
		}

		registry := microservices.NewRegistry()
		kubernetes.StartNewPodWatcher(
			client,
			time.Minute,
			"tenant,application,environment,microservice,!infrastructure",
			&microservices.Updater{
				Registry: registry,
			},
			cmd.Context().Done(),
		)

		router := mux.NewRouter()
		admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry)
		microservices.AddProxy(router, registry)
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
	viper.AutomaticEnv()
	viper.SetDefault("listenOn", "localhost:8080")
}
