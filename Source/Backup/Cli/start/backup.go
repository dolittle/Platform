package start

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongodump"
)

type mongoDump struct {
	*mongodump.MongoDump
	DumpFilename string
	DumpFilepath string
}

func CreateMongoDump(mongoHost string, dumpDir string) (*mongoDump, error) {
	dumpFileName := fmt.Sprintf("%s.gz.mongodump", time.Now().Format("2006-01-02_15-04-05"))
	dumpFilepath := fmt.Sprintf("%s/%s", dumpDir, dumpFileName)
	if !filepath.IsAbs(dumpFilepath) {
		return nil, fmt.Errorf("%s is not an absolute file path", dumpFilepath)
	}

	dump, err := createDump(mongoHost, dumpFilepath)
	if err != nil {
		return nil, err
	}

	err = dump.Init()
	if err != nil {
		return nil, err
	}

	return &mongoDump{
		MongoDump:    dump,
		DumpFilename: dumpFileName,
		DumpFilepath: dumpFilepath,
	}, nil
}

func createDump(mongoHost string, dumpFilepath string) (*mongodump.MongoDump, error) {
	toolOptions := options.New("mongodump", "", "", "", true, options.EnabledOptions{URI: true})
	_, err := toolOptions.ParseArgs([]string{mongoHost})
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
