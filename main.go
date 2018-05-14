package main

import (
	"io"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/segment-sources/mongodb/lib"
	"github.com/segmentio/objects-go"

	"github.com/segmentio/ecs-logs-go/logrus"
	"github.com/tj/docopt"
)

const (
	Version = "v0.1.5-beta"
	Usage   = `
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
)

func main() {
	m, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		logrus.Fatal(err)
	}

	if m["--debug"].(bool) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if m["--json-log"].(bool) {
		logrus.SetFormatter(logrus_ecslogs.NewFormatter())
	}

	concurrency, err := strconv.Atoi(m["--concurrency"].(string))
	if err != nil {
		logrus.Fatal(err)
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

	// If in init mode, save list of collections to schema file. Users will then have to modify the
	// file and fill in fields they want to export to their Segment warehouse.
	fileName := m["--schema"].(string)
	if config.Init {
		mongodb.InitSchema(config, fileName)
		return
	}

	// Build Segment client and define publish function for when we scan over the collections.
	writeKey := m["--write-key"].(string)
	if writeKey == "" {
		logrus.Fatal("Write key is required when not in init mode.")
	}

	description, err := mongodb.ParseSchema(fileName)
	if err == io.EOF {
		logrus.Error("Empty schema, did you run `--init`?")
	} else if err != nil {
		logrus.Fatal("Unable to parse schema", err)
	}

	segmentClient := objects.New(writeKey)
	defer segmentClient.Close()
	setWrapperFunc := func(o *objects.Object) {
		err := segmentClient.Set(o)
		if err != nil {
			logrus.WithFields(logrus.Fields{"id": o.ID, "collection": o.Collection, "properties": o.Properties}).Warn(err)
		}
	}

	logrus.Info("[%v] Mongo source started with writeKey ", Version, writeKey)
	if err := mongodb.Run(config, description, concurrency, setWrapperFunc); err != nil {
		logrus.Error("mongodb source failed to complete", err)
		os.Exit(1)
	}
}
