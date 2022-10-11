package cmd

import (
	"github.com/dolittle/platform-router/admin"
	"github.com/dolittle/platform-router/kubernetes"
	"github.com/dolittle/platform-router/microservices"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "router",
	Short: "Router",
	Long:  `Entry router`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: cmd.OutOrStdout()})

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
			microservices.Updater{
				Registry: registry,
			},
			cmd.Context().Done(),
		)

		router := mux.NewRouter()
		admin.AddApi(router.PathPrefix("/admin").Subrouter(), registry)
		//router.Handle("/admin", admin.NewAdminApi(registry))

		////router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		////	log.Print("Handling request")
		////	fmt.Fprint(w, "Hello")
		////})
		//proxy := &httputil.ReverseProxy{
		//	Director: func(r *http.Request) {
		//		vars := mux.Vars(r)
		//		log.Println(r.URL.Scheme)
		//		log.Println(vars["applicationId"], vars["environment"], vars["microserviceId"])

		//		r.URL.Scheme = "http"
		//		r.URL.Host = "127.0.0.1:7681"
		//	},
		//}

		//router := mux.NewRouter()
		//router.Handle("/{applicationId}/{environment}/{microserviceId}", proxy)

		////url, _ := url.Parse("http://localhost:7681")
		////proxy := httputil.NewSingleHostReverseProxy(url)
		server := &http.Server{
			Handler:      router,
			Addr:         viper.GetString("listenOn"),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		//log.Printf("Listening on %v", srv.Addr)
		go server.ListenAndServe()

		//log.Fatal(srv.ListenAndServe())
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

	// viper.SetDefault("tools.server.platformEnvironment", "dev")
	// viper.BindEnv("tools.server.platformEnvironment", "PLATFORM_ENVIRONMENT")
	// rootCmd.PersistentFlags().String("platform-environment", viper.GetString("tools.server.platformEnvironment"), "Platform environment (dev or prod), not linked to application environment")
	// viper.BindPFlag("tools.server.platformEnvironment", rootCmd.PersistentFlags().Lookup("platform-environment"))

	// viper.SetDefault("tools.jobs.image.operations", "dolittle/platform-operations:latest")

	// viper.BindEnv("tools.jobs.image.operations", "JOBS_OPERATIONS_IMAGE")

}
