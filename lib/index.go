package mongodb

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/segmentio/objects-go"
	"github.com/tj/go-sync/semaphore"
)

type SetObjectFunc func(o *objects.Object)

func InitSchema(config *Config, fileName string) {
	logrus.Info("Will output schema to ", fileName)
	schemaFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer schemaFile.Close()

	app := &MongoDB{}
	defer app.Close()

	logrus.Infof("Will connect to database %v@%v:%v/%v",
		config.Username, config.Hostname, config.Port, config.Database)
	// Initialize DB connection.
	if err := app.Init(config); err != nil {
		logrus.WithError(err).Error("Failed to get initialize mongo")
		return
	}

	description, err := app.GetDescription()
	if err != nil {
		logrus.WithError(err).Error("Failed to get mongo db description")
		return
	}

	if err := description.Save(schemaFile); err != nil {
		logrus.WithError(err).WithField("schema_file", schemaFile).Error("Failed to save schema file")
		return
	}

	schemaFile.Sync()
	logrus.Infof("Saved to `%s`", schemaFile.Name())
}

func ParseSchema(fileName string) (*Description, error) {
	// We must not be in init mode at this point, begin uploading source data.
	schemaFile, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer schemaFile.Close()

	return NewDescriptionFromReader(schemaFile)
}

func Run(config *Config, description *Description, concurrency int, setObjectFunc SetObjectFunc) error {
	app := &MongoDB{}
	defer app.Close()

	logrus.Infof("Will connect to database %v@%v:%v/%v",
		config.Username, config.Hostname, config.Port, config.Database)
	// Initialize DB connection.
	if err := app.Init(config); err != nil {
		logrus.Error(err)
		return err
	}

	// Launch goroutines to scan the documents in each collection.
	sem := make(semaphore.Semaphore, concurrency)

	for collection := range description.Iter() {
		// Skip collection if no fields specified in schema JSON.
		if len(collection.Fields) == 0 {
			continue
		}

		sem.Acquire()
		go func(collection *Collection, dbName string) {
			defer sem.Release()
			logrus.WithFields(logrus.Fields{"db": dbName, "collection": collection.CollectionName}).Info("Scan started")
			if err := app.ScanCollection(collection, setObjectFunc); err != nil {
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
	return nil
}
