package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/segment-sources/mongodb/lib"
	"github.com/segmentio/objects-go"
)

func main() {

	args := mongodb.ParseArgs()
	if args == nil {
		return
	}
	segmentClient := objects.New(args.WriteKey)
	defer segmentClient.Close()
	setWrapperFunc := func(o *objects.Object) {
		err := segmentClient.Set(o)
		if err != nil {
			logrus.WithFields(logrus.Fields{"id": o.ID, "collection": o.Collection, "properties": o.Properties}).Warn(err)
		}
	}
	mongodb.Run(*args, setWrapperFunc)
}
