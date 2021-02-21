package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	host        string
	port        int
	env         string
	application string
	tenant      string
	rootCmd     = &cobra.Command{
		Use:   "backups-cli",
		Short: "Perform backup",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&host, "host", "", "The host of the Backups microservice")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 80, "The port of the Backups microservice (defaults to 80)")
	rootCmd.PersistentFlags().StringVar(&application, "application", "", "The application id")
	rootCmd.PersistentFlags().StringVar(&env, "environment", "", "The environment")
	rootCmd.PersistentFlags().StringVar(&tenant, "tenant", "", "The customer tenant id")

	rootCmd.MarkPersistentFlagRequired("host")
	rootCmd.MarkPersistentFlagRequired("application")
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
