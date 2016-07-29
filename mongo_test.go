package main

import (
	"testing"

	"github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Credentials to local MongoDB instance that will be used by the tests. Make sure your local
// MongoDB has the appropriate user permissions set for this user for the test DB.
const hostname string = "localhost"
const port string = "27017"
const username string = ""
const password string = ""
const database string = "test"
const collection string = "products"

var testMongoConfig = &Config{
	Hostname: hostname,
	Port:     port,
	Username: username,
	Password: password,
	Database: database,
}

type MongoTestSuite struct {
	suite.Suite
}

func (s *MongoTestSuite) SetupTest() {
	// Try to connect to the local Mongo instance and reset relevant databases.
	session, err := connectToLocalMongo()
	if err != nil {
		return
	}
	defer session.Close()

	db := session.DB(database)
	db.DropDatabase()

	// Insert some test data.
	db.C(collection).Insert(map[string]interface{}{
		"name": "Apple",
		"cost": 1.27,
		"tags": []string{"fruit", "red"},
		"translations": map[string]interface{}{
			"spanish": "manzana",
			"french":  "pomme",
		},
	})
	db.C(collection).Insert(map[string]interface{}{
		"name": "Pear",
		"cost": 2.01,
		"tags": []string{"fruit", "yellow"},
		"translations": map[string]interface{}{
			"spanish": "pera",
		},
	})
}

func (s *MongoTestSuite) TearDownTest() {
	// Drop our test database.
	session, err := connectToLocalMongo()
	if err != nil {
		return
	}

	session.DB(database).DropDatabase()
}

func (s *MongoTestSuite) TestGetDescription() {
	t := s.T()

	app := MongoDB{}
	defer app.Close()

	// Skip test if cannot connect to local Mongo instance, but output the issue if there is one.
	err := app.Init(testMongoConfig)
	if err != nil {
		t.Skipf("Test is skipped because could not connect successfully to local Mongo instance: %v", err)
	}

	// Check to make sure description includes all of the collection names of the current DB.
	session, err := connectToLocalMongo()
	if err != nil {
		t.Fatal(err)
	}

	// Add collection names to a set.
	cNamesSet := mapset.NewSet()
	cNames, err := session.DB(database).CollectionNames()
	if err != nil {
		t.Fatal(err)
	}
	for _, cName := range cNames {
		cNamesSet.Add(cName)
	}

	desc, err := app.GetDescription()
	if err != nil {
		t.Fatal(err)
	}

	// Check that our set and desc are exactly the same set of strings.
	collectionsMap := desc.schemas[database]
	assert.Equal(t, cNamesSet.Cardinality(), len(collectionsMap), "num of collections scanned should be equal")
	for cName := range collectionsMap {
		if !cNamesSet.Contains(cName) {
			t.Fatalf("Expected MongoDB collections list to include: %v "+
				"but it didn't. Set contents are: %v", cName, cNamesSet)
		}
		cNamesSet.Remove(cName)
	}
	assert.Equal(t, 0, cNamesSet.Cardinality())
}

func (s *MongoTestSuite) TestGetForNestedKey() {
	testDoc := map[string]interface{}{
		"name": "Apple",
	}
	assert.Equal(s.T(), "Apple", getForNestedKey(testDoc, "name"))
}

func (s *MongoTestSuite) TestGetForNestedKeyNone() {
	testDoc := map[string]interface{}{
		"name": "Apple",
	}
	assert.Equal(s.T(), nil, getForNestedKey(testDoc, "nonexistentKey"))
}

func (s *MongoTestSuite) TestGetForNestedKeyNested() {
	testDoc := map[string]interface{}{
		"apple": map[string]interface{}{
			"translations": map[string]interface{}{
				"spanish": "manzana",
				"french":  "pomme",
			},
		},
	}
	assert.Equal(s.T(), "manzana", getForNestedKey(testDoc, "apple.translations.spanish"))
}

func (s *MongoTestSuite) TestGetForNestedKeyNoneNested() {
	testDoc := map[string]interface{}{
		"apple": map[string]interface{}{
			"translations": map[string]interface{}{},
		},
	}
	assert.Equal(s.T(), nil, getForNestedKey(testDoc, "apple.translations.spanish"))
}

func (s *MongoTestSuite) TestGetIdFromResultString() {
	result := map[string]interface{}{
		"_id": "abc123",
	}

	id, err := getIdFromResult(result)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), "abc123", id)
}

func (s *MongoTestSuite) TestGetIdFromResultObjectId() {
	result := map[string]interface{}{
		"_id": bson.ObjectIdHex("57881f9ce8414cf291b44b4e"),
	}

	id, err := getIdFromResult(result)
	if err != nil {
		s.T().Fatal(err)
	}

	assert.Equal(s.T(), "57881f9ce8414cf291b44b4e", id)
}

func (s *MongoTestSuite) TestGetIdFromResultNone() {
	result := map[string]interface{}{}

	_, err := getIdFromResult(result)
	assert.NotNil(s.T(), err)
}

func connectToLocalMongo() (*mgo.Session, error) {
	return mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{hostname + ":" + port},
		Database: database,
		Username: username,
		Password: password,
	})
}

func TestMongoTestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
