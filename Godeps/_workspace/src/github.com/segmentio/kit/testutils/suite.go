package testutils

import (
	"os"
	"strings"
	"testing"

	"github.com/segmentio/kit"
	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/schema"
	"github.com/stretchr/testify/suite"
)

// KitTestSuite is a testify suite augumented with some
// helper methods to make unit testing easier and straightforward
type KitTestSuite struct {
	suite.Suite
}

func (k *KitTestSuite) TearDownSuite() {
}

func (k *KitTestSuite) Run(service schema.Service) {
	// Use only the Environment configuration provider by default
	config.SetProviders([]config.ProviderType{config.Environment})

	// Because we are in the test suite here
	// we make it so that the Handler is optional
	if service.Handler == nil {
		service.Handler = func() {}
	}

	// Use a channel here because we don't want to block indefinitely
	// using kit.Run
	done := make(chan bool)
	// Run kit in a goroutine
	go kit.Run(service, done)
	// Waits for the init step to finish and return
	<-done
}

// ReplaceConfigKey replaces the configuration key with the value.
// The value must be a string because we use ENV to set the variables.
func (k *KitTestSuite) ReplaceConfigKey(key string, value string) {
	os.Setenv(strings.ToUpper(strings.Replace(key, ".", "_", -1)), value)
	config.Reload()
}

// Run runs the suite
func Run(t *testing.T, s suite.TestingSuite) {
	suite.Run(t, s)
}
