package main

import (
	"io"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/segmentio/objects-go"
	"github.com/tj/docopt"
	"github.com/tj/go-sync/semaphore"
)

const (
	Version = "0.0.1-beta"
)

var usage = `
Usage:
  mongodb
    [--debug]
    [--init]
    [--concurrency=<c>]
    [--schema=<schema-path>]
    --write-key=<segment-write-key>
    --hostname=<hostname>
    --port=<port>
    --username=<username>
    --password=<password>
    --database=<database>
  mongodb -h | --help
  mongodb --version

Options:
    "github.com/segmentio/source-db-lib/internal/domain"
  -h --help                   Show this screen
  --version                   Show version
  --write-key=<key>           Segment source write key
  --concurrency=<c>           Number of concurrent table scans [default: 1]
  --hostname=<hostname>       Database instance hostname
  --port=<port>               Database instance port number
  --username=<username>       Database instance username
  --password=<password>       Database instance password
  --database=<database>       Database instance name
  --schema=<schema-path>	    The path to the schema json file [default: schema.json]

`

func main() {
	app := &MongoDB{};
	defer app.Close()

	m, err := docopt.Parse(usage, nil, true, Version, false)
	if err != nil {
		logrus.Error(err)
		return
	}

	segmentClient := objects.New(m["--write-key"].(string))
	defer segmentClient.Close()

	setWrapperFunc := func(o *objects.Object) {
		err := segmentClient.Set(o)
		if err != nil {
			logrus.WithFields(logrus.Fields{"id": o.ID, "collection": o.Collection, "properties": o.Properties}).Warn(err)
		}
	}

	config := &Config{
		Init:         m["--init"].(bool),
		Hostname:     m["--hostname"].(string),
		Port:         m["--port"].(string),
		Username:     m["--username"].(string),
		Password:     m["--password"].(string),
		Database:     m["--database"].(string),
	}

	if m["--debug"].(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	concurrency, err := strconv.Atoi(m["--concurrency"].(string))
	if err != nil {
		logrus.Error(err)
		return
	}

	// Validate the configuration
	if _, err := govalidator.ValidateStruct(config); err != nil {
		logrus.Error(err)
		return
	}

	// Open the schema
	fileName := m["--schema"].(string)

	// Initialize DB connection.
	if err := app.Init(config); err != nil {
		logrus.Error(err)
		return
	}

	// If in init mode, save list of collections to schema file. Users will then have to modify the
	// file and fill in fields they want to export to their Segment warehouse.
	if config.Init {
		schemaFile, err := os.OpenFile(fileName, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0644)
		if err != nil {
			logrus.Error(err)
			return
		}
		defer schemaFile.Close()

		description, err := app.GetDescription()
		if err != nil {
			logrus.Error(err)
			return
		}

		if err := description.Save(schemaFile); err != nil {
			logrus.Error(err)
			return
		}

		schemaFile.Sync()
		logrus.Infof("Saved to `%s`", schemaFile.Name())
		return
	}

	// We must *not* be in init mode.
	schemaFile, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer schemaFile.Close()

	description, err := NewDescriptionFromReader(schemaFile)
	if err == io.EOF {
		logrus.Error("Empty schema, did you run `--init`?")
		return
	} else if err != nil {
		logrus.Error(err)
		return
	}

	// Launch goroutines to scan the documents in each collection.
	sem := make(semaphore.Semaphore, concurrency)

	for collection := range description.Iter() {
		// Skip collection if no fields specified.
		if len(collection.Fields) == 0 {
			continue
		}

		sem.Acquire()
		go func(collection *Collection, dbName string) {
			defer sem.Release()
			logrus.WithFields(logrus.Fields{"db": dbName, "collection": collection.CollectionName}).Info("Scan started")
			if err := app.ScanCollection(collection, setWrapperFunc); err != nil {
				logrus.Error(err)
			}
			logrus.WithFields(logrus.Fields{"db": dbName, "collection": collection.CollectionName}).Info("Scan finished")
		}(collection, app.dbName)
	}

	sem.Wait()

	// Log status
	for collection := range description.Iter() {
		logrus.WithFields(logrus.Fields{"db": app.dbName, "collection": collection.CollectionName}).Info("Sync finished")
	}
}
