package cmd

import (
	"log"

	"dolittle.io/platform/backup-cli/start"
	"github.com/spf13/cobra"
)

var (
	mongoConnectionString string
	dumpDir               string
	shareName             string
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
		backup, err := start.CreateBackup(
			backend,
			mongoDump.BackupFileName,
			tenant, env,
			eventSource,
			application,
			applicationName,
			shareName)
		if err != nil {
			return err
		}

		err = backup.NotifyStart()
		if err != nil {
			return err
		}

		log.Printf("Dumping mongo database to %s \n", mongoDump.BackupFileName)
		err = mongoDump.Dump()
		if err != nil {
			return err
		}

		err = backup.NotifyStored()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	startCmd.PersistentFlags().StringVar(&mongoConnectionString, "mongo-connection-string", "", "The mongo connection string to dump from")
	startCmd.PersistentFlags().StringVar(&dumpDir, "dump-dir", "", "The directory to put the mongo dump")
	startCmd.PersistentFlags().StringVar(&dumpDir, "share-name", "", "The share name of the azure file storage")

	startCmd.MarkPersistentFlagRequired("mongo-host")
	startCmd.MarkPersistentFlagRequired("dump-dir")
	startCmd.MarkPersistentFlagRequired("share-name")
}
