package mongodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/segmentio/go-snakecase"
	"github.com/segmentio/objects-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	db     *mgo.Database
	DBName string
}

func (m *MongoDB) Init(c *Config) error {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{c.Hostname + ":" + c.Port},
		Database: c.Database,
		Username: c.Username,
		Password: c.Password,
		Timeout:  time.Duration(5 * time.Second),
	})
	if err != nil {
		return err
	}

	m.db = session.DB(c.Database)
	m.DBName = c.Database
	logrus.Debugf("Connection to database '%s' established!", c.Database)
	return nil
}

func (m *MongoDB) GetDescription() (*Description, error) {
	desc := NewDescription()

	names, err := m.db.CollectionNames()
	if err != nil {
		return nil, err
	}

	for _, name := range names {
		// Add collections to result (it is intentionally empty right now so user can fill them out after init stage).
		desc.AddCollection(name, m.DBName)
	}

	return desc, nil
}

func (m *MongoDB) ScanCollection(c *Collection, publish func(o *objects.Object)) error {
	fieldsToInclude := make(map[string]interface{})
	for source := range c.Fields {
		fieldsToInclude[source] = 1
	}
	logrus.WithFields(logrus.Fields{"fieldsToInclude": fieldsToInclude}).Debug("Calculating which fields to include or exclude.")

	// Iterate through collection, grabbing only user specified fields.
	iter := m.db.C(c.CollectionName).Find(nil).Select(fieldsToInclude).Iter()
	var result map[string]interface{}
	for iter.Next(&result) {
		logrus.WithFields(logrus.Fields{
			"result":     result,
			"Collection": c.CollectionName,
		}).Debug("Processing row from DB")

		id, err := getIdFromResult(result)
		if err != nil {
			return err
		}

		// The destination name (e.g. name of the collection in the warehouse) can be set by the user,
		// otherwise it just defaults to the collection name in Mongo.
		var destinationName string
		if c.DestinationName == "" {
			destinationName = snakecase.Snakecase(fmt.Sprintf("%s_%s", m.DBName, c.CollectionName))
		} else {
			destinationName = c.DestinationName
		}

		// Create properties map and fill it in with all the fields were able to find.
		properties := getPropertiesMapFromResult(result, c)

		publish(&objects.Object{
			ID:         id,
			Collection: destinationName,
			Properties: properties,
		})
		logrus.WithFields(logrus.Fields{"ID": id, "Collection": destinationName, "Properties": properties}).Debug("Published row")
	}

	return iter.Close()
}

func (m *MongoDB) Close() {
	if m.db != nil {
		m.db.Session.Close()
	}
}

func getIdFromResult(result map[string]interface{}) (string, error) {
	// Translate ID from "_id" field, which can actually be one of several types.
	var id string

	switch _id := result["_id"].(type) {
	case string:
		id = _id
	case bson.ObjectId:
		id = _id.Hex()
	default:
		return "", errors.New(fmt.Sprintf("'_id' value is of unexpected type %T", result["_id"]))
	}

	return id, nil
}

func getPropertiesMapFromResult(result map[string]interface{}, c *Collection) map[string]interface{} {
	properties := make(map[string]interface{})
	for source, destination := range c.Fields {
		value := getForNestedKey(result, source)

		// The destination name (e.g. name of the field in the warehouse) can be set by the user,
		// otherwise it just defaults to the field name in Mongo.
		if destination == "" {
			destination = source
		}

		if _, ok := value.([]interface{}); ok {
			arrayJSON, err := json.Marshal(value)
			if err != nil {
				logrus.Errorf("[Error] Unable to marshall value. Skipping `%v` err: %v", value, err)
			} else {
				properties[destination] = string(arrayJSON)
			}
		} else if value != nil && value != bson.Undefined {
			properties[destination] = value
		}
	}
	return properties
}

// Searches for a value in the map if the key (which may refer to a nested field several levels deep).
// If that value cannot be found, returns nil. For example, if the key "inner_dict.key_1" is passed in,
// this method looks for a dict called inner_dict and then for a field keyed by "key_1" in that dict.
func getForNestedKey(curMap map[string]interface{}, key string) interface{} {
	if curMap == nil {
		return nil
	}

	firstDot := strings.Index(key, ".")
	if firstDot == -1 {
		return curMap[key]
	}

	curKey, nextKey := key[:firstDot], key[firstDot+1:]
	if val, ok := curMap[curKey]; ok {
		if val, ok := val.(map[string]interface{}); ok {
			return getForNestedKey(val, nextKey)
		}
	}
	return nil
}
