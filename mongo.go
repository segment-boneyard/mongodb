package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/segmentio/go-source"
	"github.com/segmentio/kit/log"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func syncMongo(ctx context.Context, session *mgo.Session, sourceClient source.Client) {
	databases, err := session.DatabaseNames()
	check(err)

	var wg sync.WaitGroup
	for _, database := range databases {
		wg.Add(1)
		go func(database string) {
			defer wg.Done()
			db := session.DB(database)
			syncDatabase(context.WithValue(ctx, "database", database), db, sourceClient)
		}(database)
	}
	wg.Wait()
}

func syncDatabase(ctx context.Context, db *mgo.Database, sourceClient source.Client) {
	collections, err := db.CollectionNames()
	check(err)

	var wg sync.WaitGroup
	for _, collection := range collections {
		wg.Add(1)
		go func(collection string) {
			defer wg.Done()
			c := db.C(collection)
			syncCollection(context.WithValue(ctx, "collection", collection), c, sourceClient)
		}(collection)
	}
	wg.Wait()
}

func syncCollection(ctx context.Context, collection *mgo.Collection, sourceClient source.Client) {
	log.With(map[string]interface{}{
		"database":   ctx.Value("database"),
		"collection": ctx.Value("collection"),
	}).Infof("syncing collection")

	n := 0
	iter := collection.Find(nil).Snapshot().Iter()
	var elem bson.M
	for ; iter.Next(&elem); n++ {
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

	log.With(map[string]interface{}{
		"database":   ctx.Value("database"),
		"collection": ctx.Value("collection"),
		"count":      n,
	}).Infof("synced collection")
}
