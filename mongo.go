package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/segmentio/go-source"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func syncMongo(session *mgo.Session, sourceClient source.Client) {
	names, err := session.DatabaseNames()
	check(err)

	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			db := session.DB(name)
			syncDatabase(db, sourceClient)
		}(name)
	}
	wg.Wait()
}

func syncDatabase(db *mgo.Database, sourceClient source.Client) {
	collections, err := db.CollectionNames()
	check(err)

	var wg sync.WaitGroup
	for _, collection := range collections {
		wg.Add(1)
		go func(collection string) {
			defer wg.Done()
			c := db.C(collection)
			syncCollection(c, sourceClient)
		}(collection)
	}
	wg.Wait()
}

func syncCollection(collection *mgo.Collection, sourceClient source.Client) {
	iter := collection.Find(nil).Snapshot().Iter()
	var elem bson.M
	for iter.Next(&elem) {
		var id string
		switch _id := elem["_id"].(type) {
		case string:
			id = _id
		case bson.ObjectId:
			id = _id.String()
		default:
			panic(errors.New(fmt.Sprintf("unknown type for _id: %T", elem["_id"])))
		}
		delete(elem, "_id")
		err := sourceClient.Set(collection.Name, id, elem)
		check(err)
	}
	check(iter.Close())
}
