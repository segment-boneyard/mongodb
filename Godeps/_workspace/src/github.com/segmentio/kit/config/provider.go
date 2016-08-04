package config

import "github.com/segmentio/kit/schema"

// ProviderType defines a configuration provider
// as an alias for int, values are defined as const
// in this file. Users may use these values to define priorities.
type ProviderType int

const (
	CommandLine ProviderType = iota
	Environment
)

// ConfigProvider is the interface that every configuration provider
// must satisfy in order to be included.
type ConfigProvider interface {
	Setup(service schema.Service) error
	Get(value schema.ConfigValue) interface{}
}
