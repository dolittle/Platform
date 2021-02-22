package cmd

import (
	"log"

	"dolittle.io/system-data/backups-cli/backups"
	"dolittle.io/system-data/backups-cli/mongo"
	"github.com/spf13/cobra"
)

var (
	mongoConnectionString string
	archive               string
	shareName             string
)

func init() {
	startCmd.PersistentFlags().StringVar(&archive, "archive", "", "The directory to put the mongo dump")
	startCmd.PersistentFlags().StringVar(&shareName, "share-name", "", "The share name of the azure file storage")

	startCmd.MarkPersistentFlagRequired("archive")
	startCmd.MarkPersistentFlagRequired("share-name")
}

var startCmd = &cobra.Command{
	Use:   "mongo [connection string]",
	Short: "Does a mongodump",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mongoConnectionString = args[0]
		log.Println("Starting backup routine")
		log.Printf("Setting up system for committing event to Backup microservice at endpoint %s:%d\n", backendURL, backendPort)
		backupFileName := mongo.CreateBackupFileName()
		backups, err := backups.CreateBackups(
			backendURL,
			backendPort,
			tenant,
			environment,
			application,
			shareName)
		if err != nil {
			return err
		}
		log.Printf("Initialzing mongo dump with connection string %s", mongoConnectionString)
		mongoDump, err := mongo.CreateMongoDump(mongoConnectionString, archive, backupFileName)
		if err != nil {
			return err
		}

		log.Printf("Dumping mongo database to %s \n", mongoDump.DumpFilePath)
		backupDurationInSeconds, err := mongoDump.Dump()
		if err != nil {
			return err
		}

		err = backups.NotifyStored(backupFileName, backupDurationInSeconds)
		if err != nil {
			return err
		}

		return nil
	},
}

// func handleDatabaseError(err error, backups *backups.Backups, backupFileName string) error {
// 	failureReason := err.Error()
// 	log.Printf("%s\n", failureReason)
// 	err = backups.NotifyFailed(backupFileName, failureReason)
// 	if err != nil {
// 		log.Println("An error occurred while notifying of failed backup")
// 		return err
// 	}
// 	return err
// }
