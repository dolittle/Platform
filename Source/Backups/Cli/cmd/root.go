package cmd

import (
	"fmt"
	"os"

	"dolittle.io/system-data/backups-cli/config"
	"github.com/spf13/cobra"
)

var (
	backendURL  string
	backendPort int
	application string
	tenant      string
	environment string
	rootCmd     = &cobra.Command{
		Use:   "backups-cli",
		Short: "Perform backup",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Initialize(cmd)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&backendURL, "backend-url", "", "The url of the Backups microservice")
	rootCmd.MarkPersistentFlagRequired("backend-url")
	rootCmd.PersistentFlags().IntVarP(&backendPort, "backend-port", "p", 80, "The port of the Backups microservice (defaults to 80)")
	rootCmd.PersistentFlags().StringVarP(&application, "application", "a", "", "The application id")
	rootCmd.MarkPersistentFlagRequired("application")
	rootCmd.PersistentFlags().StringVarP(&tenant, "tenant", "t", "", "The customer tenant id")
	rootCmd.MarkPersistentFlagRequired("tenant")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "The environment")
	rootCmd.MarkPersistentFlagRequired("environment")

	rootCmd.AddCommand(startCmd)
}

// Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
