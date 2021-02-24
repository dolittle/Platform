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
		backups, mongoDump, err := createDependencies()
		if err != nil {
			return err
		}

		log.Printf("Dumping mongo database to %s \n", mongoDump.DumpFilePath)
		backupDurationInSeconds, err := mongoDump.Dump()
		if err != nil {
			return err
		}

		if err = backups.NotifyStored(mongoDump.DumpFileName, backupDurationInSeconds); err != nil {
			log.Printf("Failed to notify Backups microservice of persisted backup: %s", err.Error())
		}

		return nil
	},
}

func createDependencies() (*backups.Backups, *mongo.MongoDump, error) {
	log.Printf("Setting up system for committing event to Backup microservice at endpoint %s:%d\n", backendURL, backendPort)
	backups, err := backups.CreateBackups(
		backendURL,
		backendPort,
		tenant,
		environment,
		application,
		shareName)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Initialzing mongo dump with connection string %s", mongoConnectionString)
	mongoDump, err := mongo.CreateMongoDump(mongoConnectionString, archive)
	if err != nil {
		return nil, nil, err
	}
	return backups, mongoDump, nil
}
