package main

import (
	"testing"

	"camlistore.org/pkg/test/dockertest"

	"github.com/bmizerany/assert"
	"github.com/segmentio/go-source/mock"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

func TestReal(t *testing.T) {
	t.Skip()

	session, err := mgo.Dial("localhost")
	check(err)

	source := mock.New()
	go func() {
		for {
			<-source.SetCalls
		}
	}()

	syncMongo(context.Background(), "test", session, source)
}

func TestSync(t *testing.T) {
	container, ip := dockertest.SetupMongoContainer(t)
	defer container.KillRemove(t)

	session, err := mgo.Dial(ip)
	assert.Equal(t, nil, err)
	collection := session.DB("testdb").C("testcollection")
	data := map[string]interface{}{
		"_id":  "foo",
		"name": "Pratek",
	}
	if err := collection.Insert(data); err != nil {
		t.Error(err)
	}

	source := mock.New()
	assert.Equal(t, nil, err)

	syncMongo(context.Background(), "testdb", session, source)

	c := <-source.SetCalls
	assert.Equal(t, c.ID, "foo") // todo: test to verify always a string not hex.
	assert.Equal(t, c.Collection, "testcollection")
	assert.Equal(t, c.Properties, map[string]interface{}{
		"name": "Pratek",
	})
}
