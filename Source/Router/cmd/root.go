package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "platform-router",
	Short: "platform-router is a router for Dolittle microservice ports in a Kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

// Execute starts the cobra.Command CLI execution
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringSlice("config", nil, "A configuration file to load, can be specified multiple times")
	rootCmd.PersistentFlags().String("logger.format", "console", "The logging format to use, 'json' or 'console'")
	rootCmd.PersistentFlags().String("logger.level", "info", "The logging minimum log level to output")

	rootCmd.AddCommand(proxyCmd)
}
