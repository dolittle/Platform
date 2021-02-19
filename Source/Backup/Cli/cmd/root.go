package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	backend         string
	eventSource     string
	env             string
	application     string
	applicationName string
	tenant          string
	rootCmd         = &cobra.Command{
		Use:   "platform-backup",
		Short: "Initiate backups",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&backend, "backend", "", "The backend uri of the Backup Microservice")
	rootCmd.PersistentFlags().StringVar(&eventSource, "event-source", "", "The event source id of the backup (unique for each environment and tenant and application)")
	rootCmd.PersistentFlags().StringVar(&application, "application", "", "The application id")
	rootCmd.PersistentFlags().StringVar(&applicationName, "application-name", "", "The application name")
	rootCmd.PersistentFlags().StringVar(&env, "environment", "", "The environment")
	rootCmd.PersistentFlags().StringVar(&tenant, "tenant", "", "The customer tenant id")

	rootCmd.MarkPersistentFlagRequired("backend")
	rootCmd.MarkPersistentFlagRequired("event-source")
	rootCmd.MarkPersistentFlagRequired("application")
	rootCmd.MarkPersistentFlagRequired("application-name")
	rootCmd.MarkPersistentFlagRequired("environment")
	rootCmd.MarkPersistentFlagRequired("tenant")
	rootCmd.AddCommand(startCmd)
}

// Execute starts the program
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
