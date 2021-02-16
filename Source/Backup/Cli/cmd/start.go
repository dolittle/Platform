package cmd

import (
	"log"

	"dolittle.io/platform/backup-cli/start"
	"github.com/spf13/cobra"
)

type startBackupPayload struct {
	DumpFilename string `json:"dumpFilename"`
	Tenant       string `json:"tenant"`
	Environment  string `json:"environment"`
	EventSource  string `json:"eventSource"`
	Application  string `json:"application"`
}

var (
	mongoConnectionString string
	dumpDir               string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Starting backup routine")
		log.Printf("Initialzing mongo dump with connection string %s dumping to %s", mongoConnectionString, dumpDir)
		mongoDump, err := start.CreateMongoDump(mongoConnectionString, dumpDir)
		if err != nil {
			return err
		}

		log.Printf("Setting up system for committing event to Backup microservice at endpoint %s\n", backend)
		backupStarted, err := start.CreateBackupStarted(backend)
		if err != nil {
			return err
		}
		log.Printf("Dumping mongo database to %s \n", mongoDump.DumpFilepath)
		err = mongoDump.Dump()
		if err != nil {
			return err
		}

		log.Println("Committing event")
		err = backupStarted.Commit(mongoDump.DumpFilename, tenant, env, eventSource, application)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	startCmd.PersistentFlags().StringVar(&mongoConnectionString, "mongo-connection-string", "", "The mongo connection string to dump from")
	startCmd.PersistentFlags().StringVar(&dumpDir, "dump-dir", "", "The directory to put the mongo dump")

	startCmd.MarkPersistentFlagRequired("mongo-host")
	startCmd.MarkPersistentFlagRequired("dump-dir")
}
