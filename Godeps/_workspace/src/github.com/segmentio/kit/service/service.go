package service

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/log"
	"github.com/segmentio/kit/schema"
	"github.com/segmentio/kit/stats"
	"github.com/segmentio/kit/stats/runtime"
	"gopkg.in/validator.v2"
)

// Init validates the serviceSchema provided, initializes configuration,
// metrics, and the bundled runtime reporter
func Init(serviceSchema schema.Service) error {

	// Validate the Service Schema
	if err := validator.Validate(serviceSchema); err != nil {
		return err
	}

	// Initialize Configuration
	// Config is configured at package level, this ensures better readability
	// and less state in different structs. Its role is to provide a consistent
	// way to get key/value objects from different providers in a merge fashion.
	// For this reason, the config is initialized globally.
	if err := config.Init(serviceSchema); err != nil {
		return err
	}

	// Run pprof by default on port 6060
	go func() {
		log.With(log.M{"host": "localhost", "port": "6060"}).Debug("starting pprof")
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	// Initialize Metrics
	if err := stats.Init(serviceSchema); err != nil {
		return err
	}

	// Initialize runtime reporter
	if err := runtime.Init(serviceSchema); err != nil {
		return err
	}

	// TODO(vince) Add healthcheck endpoint on fixed port number, maybe same as pprof?
	// TODO(vince) Resources initialization step

	// Run the handler
	if serviceSchema.Handler == nil {
		return fmt.Errorf("Handler is nil")
	}

	// Executes the service handler in a goroutine
	go serviceSchema.Handler()

	// Finish Init logging and returning
	log.With(log.M{"name": serviceSchema.Name, "version": serviceSchema.Version}).Infof("Service started")
	return nil
}
