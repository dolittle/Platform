package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "router",
	Short: "Router",
	Long:  `Entry router`,
}

func Execute() {
	log.Print("Starting server")
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Handling request")
		fmt.Fprint(w, "Hello")
	})

	srv := &http.Server{
		Handler:      router,
		Addr:         viper.GetString("listenOn"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Listening on %v", srv.Addr)

	log.Fatal(srv.ListenAndServe())
}

func init() {
	log.Print("Performing init")
	viper.AutomaticEnv()

	viper.SetDefault("listenOn", "localhost:8080")

	// rootCmd.AddCommand(api.RootCmd)

	// viper.SetDefault("tools.server.platformEnvironment", "dev")
	// viper.BindEnv("tools.server.platformEnvironment", "PLATFORM_ENVIRONMENT")
	// rootCmd.PersistentFlags().String("platform-environment", viper.GetString("tools.server.platformEnvironment"), "Platform environment (dev or prod), not linked to application environment")
	// viper.BindPFlag("tools.server.platformEnvironment", rootCmd.PersistentFlags().Lookup("platform-environment"))

	// viper.SetDefault("tools.jobs.image.operations", "dolittle/platform-operations:latest")

	// viper.BindEnv("tools.jobs.image.operations", "JOBS_OPERATIONS_IMAGE")

}
