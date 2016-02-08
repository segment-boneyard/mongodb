package config

import (
	"fmt"
	"reflect"
	"sync/atomic"

	"github.com/segmentio/kit/log"
	"github.com/segmentio/kit/schema"
)

var (
	defaultConfigPtr atomic.Value
	defaultProviders = []ConfigProvider{&docOptProvider{}, &envProvider{}}
	providedSchema   schema.Config
)

// SetProviders replaces the current default providers used
// when building the Config object.
func SetProviders(configProviders []ProviderType) {
	defaultProviders = []ConfigProvider{}
	for _, p := range configProviders {
		switch p {
		case CommandLine:
			defaultProviders = append(defaultProviders, &docOptProvider{})
		case Environment:
			defaultProviders = append(defaultProviders, &envProvider{})
		}
	}
}

// Init takes care of initializing the configuration schema and populate the values.
// Looks up the registered providers in order, the first one to have a value
// for the specified Key is used.
func Init(service schema.Service) error {
	log.Debug("[CONFIG] Initializing providers")

	if len(service.Config) == 0 {
		log.Debug("[CONFIG] Empty configuration provided, skipping")
		emptyConfig := map[string]interface{}{}
		defaultConfigPtr.Store(emptyConfig)
		return nil
	}

	var err error
	for i, p := range defaultProviders {
		providerName := reflect.ValueOf(p).Elem().Type().Name()
		log.Debugf("\t -> Priority Level: %d - %s", len(defaultProviders)-i, providerName)
		err = p.Setup(service)
		if err != nil {
			return err
		}
	}

	providedSchema = service.Config

	if err := Reload(); err != nil {
		return err
	}

	return nil
}

// Reload takes care of allocating a new configuration map[string]interface{},
// query every provider registered in order and get a value for the key queried
// if no value is returned, the default one may be used (if any).
// After populating the new configuration map, it takes care of replacing the
// configuration pointer atomically, ensuring thread-safety.
func Reload() error {
	configMap := map[string]interface{}{}

	for _, val := range providedSchema {
		log.Debugf("[CONFIG] Looking a value for key `%s`", val.Key)
		found := false
		// Range over the providers to find a value
		for _, p := range defaultProviders {
			providerName := reflect.ValueOf(p).Elem().Type().Name()
			log.Debugf("\t -> Checking provider %s", providerName)
			value := p.Get(val)
			if value != nil {
				found = true
				configMap[val.Key] = value
				log.Debugf("\t -> Setting %s's value for key `%s`: `%v`", providerName, val.Key, value)
				break
			}
		}
		// If the value is not found
		if !found {
			if val.Default != nil {
				configMap[val.Key] = val.Default
				log.Debugf("\t -> Setting default value for `%s`: `%v`", val.Key, val.Default)
			} else if val.Required {
				log.Debugf("\t -> No value found for required key `%s`", val.Key)
				return fmt.Errorf("[CONFIG] Required key `%s` was not found in any provider", val.Key)
			} else {
				log.Debugf("\t -> Setting `nil` for key %s because it's optional and nothing was found", val.Key)
				configMap[val.Key] = nil
			}
		}
	}

	// Replace the pointer
	log.Debug("[CONFIG] Replacing configuration")
	defaultConfigPtr.Store(configMap)
	return nil
}

// Get returns the value as an interface{}
func Get(key string) interface{} {
	defaultConfig := defaultConfigPtr.Load().(map[string]interface{})
	if val, ok := defaultConfig[key]; ok {
		return val
	}
	log.Warnf("[CONFIG] Key `%s` is missing, have you forgot to include it in the schema?", key)
	return nil
}

// Returns an interface and a boolean, to check if the value was set
func GetOk(key string) (interface{}, bool) {
	defaultConfig := defaultConfigPtr.Load().(map[string]interface{})
	if val, ok := defaultConfig[key]; ok && val != nil {
		return val, ok
	}
	return nil, false
}
