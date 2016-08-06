package main

import (
	"io"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/segment-sources/mongodb/lib"
	"github.com/segmentio/ecs-logs-go/logrus"
	"github.com/segmentio/objects-go"
	"github.com/tj/docopt"
	"github.com/tj/go-sync/semaphore"
)

const (
	Version = "v0.1.2-beta"
)

var usage = `
Usage:
  mongodb
    [--debug]
    [--init]
    [--json-log]
    [--concurrency=<c>]
    [--schema=<schema-path>]
    [--write-key=<segment-write-key>]
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
	[--debug]										Set logrus level to .DebugLevel
	[--json-log]								Format log as JSON. Useful for ecs-logs for example
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
	app := &mongodb.MongoDB{}
	defer app.Close()

	m, err := docopt.Parse(usage, nil, true, Version, false)
	if err != nil {
		logrus.Error(err)
		return
	}

	if m["--debug"].(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if m["--json-log"].(bool) {
		logrus.SetFormatter(logrus_ecslogs.NewFormatter())
	}

	concurrency, err := strconv.Atoi(m["--concurrency"].(string))
	if err != nil {
		logrus.Error(err)
		return
	}

	// Load and validate DB configuration.
	config := &mongodb.Config{
		Init:     m["--init"].(bool),
		Hostname: m["--hostname"].(string),
		Port:     m["--port"].(string),
		Username: m["--username"].(string),
		Password: m["--password"].(string),
		Database: m["--database"].(string),
	}
	_, err = govalidator.ValidateStruct(config)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Infof("Will connect to database %v@%v:%v/%v",
		config.Username, config.Hostname, config.Port, config.Database)
	// Initialize DB connection.
	if err := app.Init(config); err != nil {
		logrus.Error(err)
		return
	}
	// If in init mode, save list of collections to schema file. Users will then have to modify the
	// file and fill in fields they want to export to their Segment warehouse.
	fileName := m["--schema"].(string)
	if config.Init {
		initSchema(fileName, app)
		return
	}

	// We must not be in init mode at this point, begin uploading source data.
	schemaFile, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer schemaFile.Close()

	description, err := mongodb.NewDescriptionFromReader(schemaFile)
	if err == io.EOF {
		logrus.Error("Empty schema, did you run `--init`?")
		return
	} else if err != nil {
		logrus.Error(err)
		return
	}

	// Build Segment client and define publish function for when we scan over the collections.
	writeKey := m["--write-key"].(string)
	if writeKey == "" {
		logrus.Error("Write key is required when not in init mode.")
		return
	}
	run(writeKey, app, description, concurrency)
}

func initSchema(fileName string, app *mongodb.MongoDB) {
	logrus.Info("Will output schema to ", fileName)
	schemaFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
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
}

func run(writeKey string, app *mongodb.MongoDB, description *mongodb.Description, concurrency int) {
	logrus.Info("Mongo source started with writeKey ", writeKey)
	segmentClient := objects.New(writeKey)
	defer segmentClient.Close()

	setWrapperFunc := func(o *objects.Object) {
		err := segmentClient.Set(o)
		if err != nil {
			logrus.WithFields(logrus.Fields{"id": o.ID, "collection": o.Collection, "properties": o.Properties}).Warn(err)
		}
	}

	// Launch goroutines to scan the documents in each collection.
	sem := make(semaphore.Semaphore, concurrency)

	for collection := range description.Iter() {
		// Skip collection if no fields specified in schema JSON.
		if len(collection.Fields) == 0 {
			continue
		}

		sem.Acquire()
		go func(collection *mongodb.Collection, dbName string) {
			defer sem.Release()
			logrus.WithFields(logrus.Fields{"db": dbName, "collection": collection.CollectionName}).Info("Scan started")
			if err := app.ScanCollection(collection, setWrapperFunc); err != nil {
				logrus.Error(err)
			}
			logrus.WithFields(logrus.Fields{"db": dbName, "collection": collection.CollectionName}).Info("Scan finished")
		}(collection, app.DBName)
	}

	sem.Wait()

	// Log status
	for collection := range description.Iter() {
		logrus.WithFields(logrus.Fields{"db": app.DBName, "collection": collection.CollectionName}).Info("Sync finished")
	}
}
