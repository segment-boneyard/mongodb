package stats

import (
	"fmt"
	"time"

	"github.com/ooyala/go-dogstatsd"
	"github.com/segmentio/kit/config"
	"github.com/segmentio/kit/schema"
)

var (
	datadogClient *dogstatsd.Client
)

// Init initializes the stats package using the serviceSchema
// by default it uses datadogClient, sets the Namespace to the service Name
// and reports every metric with program:<SERVICE_NAME> version:<SERVICE_VERSION>
// tags to Datadog.
//
// NOTE: in future we may want to abstract this out and have it pluggable with
// different providers.
func Init(serviceSchema schema.Service) error {
	var err error

	datadogAddress, ok := config.GetOk("datadog.addr")
	if !ok {
		datadogAddress = "0.0.0.0:8125"
	}

	datadogClient, err = dogstatsd.New(datadogAddress.(string))
	if err != nil {
		return err
	}

	datadogClient.Namespace = fmt.Sprintf("%s.", serviceSchema.Name)
	datadogClient.Tags = []string{
		fmt.Sprintf("program:%s", serviceSchema.Name),
		fmt.Sprintf("version:%s", serviceSchema.Version),
	}

	return nil
}

// Increment increments the counter for the given bucket.
func Increment(name string, count int, rate float64, tags ...string) error {
	return datadogClient.Count(name, int64(count), tags, rate)
}

// Incr increments the counter for the given bucket by 1 at a rate of 1.
func Incr(name string, tags ...string) error {
	return Increment(name, 1, 1, tags...)
}

// IncrBy increments the counter for the given bucket by N at a rate of 1.
func IncrBy(name string, n int, tags ...string) error {
	return Increment(name, n, 1, tags...)
}

// Decrement decrements the counter for the given bucket.
func Decrement(name string, count int, rate float64, tags ...string) error {
	return Increment(name, -count, rate, tags...)
}

// Decr decrements the counter for the given bucket by 1 at a rate of 1.
func Decr(name string, tags ...string) error {
	return Increment(name, -1, 1, tags...)
}

// DecrBy decrements the counter for the given bucket by N at a rate of 1.
func DecrBy(name string, value int, tags ...string) error {
	return Increment(name, -value, 1, tags...)
}

// Duration records time spent for the given bucket with time.Duration.
func Duration(name string, duration time.Duration, tags ...string) error {
	return Histogram(name, millisecond(duration), tags...)
}

// Histogram is an alias of .Duration() until the statsd protocol figures its shit out.
func Histogram(name string, value int, tags ...string) error {
	return datadogClient.Histogram(name, float64(value), tags, 1)
}

// Gauge records arbitrary values for the given bucket.
func Gauge(name string, value int, tags ...string) error {
	return datadogClient.Gauge(name, float64(value), tags, 1)
}

// Annotate sends an annotation.
func Annotate(name string, value string, args ...interface{}) error {
	return datadogClient.Event(name, value, nil)
}

func millisecond(d time.Duration) int {
	return int(d.Seconds() * 1000)
}
