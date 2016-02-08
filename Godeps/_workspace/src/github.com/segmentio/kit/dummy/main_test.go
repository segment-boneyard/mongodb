package main

import (
	"testing"

	"github.com/segmentio/kit/schema"
	"github.com/segmentio/kit/testutils"
)

type DummyTestSuite struct {
	testutils.KitTestSuite
}

func (d *DummyTestSuite) SetupSuite() {
	d.Run(schema.Service{
		Name:    "dummy",
		Version: "0.0.0",
		Config: schema.Config{
			schema.ConfigValue{
				Key:     "abc",
				Default: "123",
			},
		},
	})
}

func (d *DummyTestSuite) TestWhatever() {
	d.Equal("123", configAbcGetter())
	d.ReplaceConfigKey("abc", "cba")
	d.Equal("cba", configAbcGetter())
}

func TestDummy(t *testing.T) {
	testutils.Run(t, new(DummyTestSuite))
}
