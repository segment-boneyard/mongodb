package kit

import (
	"os"

	"github.com/segmentio/kit/log"
	"github.com/segmentio/kit/schema"
	"github.com/segmentio/kit/service"
)

// Run initializes and runs the service,
// it exits if there are any errors. It takes a list of optional channels,
// as varargs and it notifies each channel upon successful initialization.
func Run(serviceSchema schema.Service, done ...chan bool) {
	if err := log.Init(serviceSchema); err != nil {
		log.Errorf("Error initializing logger:  %s", err.Error())
		os.Exit(1)
	}

	if err := service.Init(serviceSchema); err != nil {
		log.Errorf("Error initializating the service: %s", err.Error())
		os.Exit(1)
	}

	if len(done) > 0 {
		for _, ch := range done {
			ch <- true
		}
	}

	// Lock forever
	select {}
}
