package mongo

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
)

type mongoDump struct {
	internal     *mongodump.MongoDump
	DumpFilePath string
}

func CreateBackupFileName() string {
	return fmt.Sprintf("%s.gz.mongodump", time.Now().Format("2006-01-02_15-04-05"))
}

func CreateMongoDump(mongoHost string, dumpDir string, backupFileName string) (*mongoDump, error) {
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

	return &mongoDump{
		internal:     dump,
		DumpFilePath: dumpFilePath,
	}, nil
}

func (m *mongoDump) Dump() error {
	err := m.internal.Dump()
	if err != nil {
		return fmt.Errorf("Failed dumping database: %s", err.Error())
	}
	return nil
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

func createInitError(errorMessage string) error {
	return fmt.Errorf("Failed initializing mongo dump tool: %s", errorMessage)
}
