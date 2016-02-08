package main

import (
	"os"

	"github.com/segmentio/go-source"
	"github.com/segmentio/kit"
	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/schema"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

var Version = "0.1.0"

func main() {
	kit.Run(schema.Service{
		Name:    "mongodb",
		Version: Version,
		Config: schema.Config{
			{
				Key:      "url",
				Required: true,
			},
		},
		Handler: run,
	})
}

func run() {
	url := config.Get("url").(string)
	session, err := mgo.Dial(url)
	check(err)

	sourceClient, err := source.New(&source.Config{
		URL: "http://localhost:4000/rpc",
	})
	check(err)

	syncMongo(context.Background(), session, sourceClient)

	os.Exit(0)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
