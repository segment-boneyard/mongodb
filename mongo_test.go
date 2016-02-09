package main

import (
	"testing"

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
