package cmd

import (
	"log"

	"dolittle.io/system/backups-cli/start"
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
		log.Printf("Setting up system for committing event to Backup microservice at endpoint %s:%d\n", host, port)
		backupFileName := start.CreateBackupFileName()
		backups, err := start.CreateBackups(
			host,
			port,
			tenant,
			env,
			application,
			shareName)
		if err != nil {
			return err
		}
		log.Printf("Initialzing mongo dump with connection string %s", mongoConnectionString)
		mongoDump, err := start.CreateMongoDump(mongoConnectionString, dumpDir, backupFileName)
		if err != nil {
			return handleDatabaseError(err, backups, backupFileName)
		}

		log.Printf("Dumping mongo database to %s \n", mongoDump.DumpFilePath)
		err = mongoDump.Dump()
		if err != nil {
			return handleDatabaseError(err, backups, backupFileName)
		}

		err = backups.NotifyStored(backupFileName)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	startCmd.PersistentFlags().StringVar(&mongoConnectionString, "mongo-connection-string", "", "The mongo connection string to dump from")
	startCmd.PersistentFlags().StringVar(&dumpDir, "dump-dir", "", "The directory to put the mongo dump")
	startCmd.PersistentFlags().StringVar(&shareName, "share-name", "", "The share name of the azure file storage")

	startCmd.MarkPersistentFlagRequired("mongo-host")
	startCmd.MarkPersistentFlagRequired("dump-dir")
	startCmd.MarkPersistentFlagRequired("share-name")
}

func handleDatabaseError(err error, backups *start.Backups, backupFileName string) error {
	failureReason := err.Error()
	log.Printf("%s\n", failureReason)
	err = backups.NotifyFailed(backupFileName, failureReason)
	if err != nil {
		log.Println("An error occurred while notifying of failed backup")
		return err
	}
	return err
}
