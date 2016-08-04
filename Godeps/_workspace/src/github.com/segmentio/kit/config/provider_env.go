package config

import (
	"os"
	"strings"

	"github.com/segmentio/kit/schema"
)

type envProvider struct{}

func (e *envProvider) Setup(service schema.Service) error {
	return nil
}

func (e *envProvider) Get(val schema.ConfigValue) interface{} {
	key := strings.ToUpper(strings.Replace(val.Key, ".", "_", -1))
	if k := os.Getenv(key); k != "" {
		return k
	}
	return nil
}
