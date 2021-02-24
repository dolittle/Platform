package mongo

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
)

type MongoDump struct {
	internal     *mongodump.MongoDump
	DumpFilePath string
	DumpFileName string
}

func CreateMongoDump(mongoHost string, dumpDir string) (*MongoDump, error) {
	backupFileName := createBackupFileName()
	dumpFilePath := fmt.Sprintf("%s/%s", dumpDir, backupFileName)
	if !filepath.IsAbs(dumpFilePath) {
		return nil, createInitError(fmt.Sprintf("Dump file path %s is not an absolute file path", dumpFilePath))
	}

	dump, err := createDump(mongoHost, dumpFilePath)
	if err != nil {
		return nil, createInitError(err.Error())
	}

	err = dump.Init()
	if err != nil {
		return nil, createInitError(err.Error())
	}

	return &MongoDump{
		internal:     dump,
		DumpFilePath: dumpFilePath,
		DumpFileName: backupFileName,
	}, nil
}

func (m *MongoDump) Dump() (backupDurationInSeconds uint, err error) {
	start := time.Now()
	err = m.internal.Dump()
	elapsed := time.Since(start)
	backupDurationInSeconds = uint(elapsed.Round(time.Second).Seconds())
	if err != nil {
		err = fmt.Errorf("Failed dumping database: %s", err.Error())
		return backupDurationInSeconds, err
	}
	log.Printf("Dumping took %s\n", elapsed.String())
	return backupDurationInSeconds, err
}

func createDump(mongoConnectionString string, dumpFilepath string) (*mongodump.MongoDump, error) {
	toolOptions := options.New("mongodump", "", "", "", true, options.EnabledOptions{Connection: true})
	_, err := toolOptions.ParseArgs([]string{mongoConnectionString})
	if err != nil {
		return nil, err
	}

	inputOptions := &mongodump.InputOptions{}
	outputOptions := &mongodump.OutputOptions{
		Gzip:                   true,
		Archive:                dumpFilepath,
		NumParallelCollections: 1,
	}
	return &mongodump.MongoDump{
		ToolOptions:   toolOptions,
		InputOptions:  inputOptions,
		OutputOptions: outputOptions,
	}, nil
}

func createBackupFileName() string {
	return fmt.Sprintf("%s.gz.mongodump", time.Now().Format("2006-01-02_15-04-05"))
}

func createInitError(errorMessage string) error {
	return fmt.Errorf("Failed initializing mongo dump tool: %s", errorMessage)
}
