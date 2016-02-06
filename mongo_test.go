package main

import (
	"testing"

	"github.com/segmentio/go-source/mock"

	"gopkg.in/mgo.v2"
)

func TestReal(t *testing.T) {
	session, err := mgo.Dial("localhost")
	check(err)

	source := mock.New()
	go func() {
		for {
			<-source.SetCalls
		}
	}()

	syncMongo(session, source)
}
