package config

import (
	"os"
	"strings"
	"testing"

	"github.com/segmentio/kit/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(c schema.Config) error {
	return Init(schema.Service{
		Name:    "kit",
		Version: "internal",
		Config:  c,
	})
}

func TestSetProvidersPriority(t *testing.T) {
	SetProviders([]ProviderType{Environment})
	os.Setenv("REQUIRED_ON_ENV", "hereiam")
	err := setupTestConfig(schema.Config{
		{
			Key:      "required.on.env",
			Required: true,
		},
	})
	assert.Len(t, defaultProviders, 1)
	assert.Nil(t, err)
	assert.Equal(t, "hereiam", Get("required.on.env").(string))
}

func TestConfigValueRequiredError(t *testing.T) {
	err := setupTestConfig(schema.Config{
		{
			Key:      "required.value.but.not.found",
			Required: true,
		},
	})
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "required.value.but.not.found"))
}

func TestConfigValueOptionalNil(t *testing.T) {
	err := setupTestConfig(schema.Config{
		{
			Key: "optional.value.but.not.found",
		},
	})
	assert.NoError(t, err)
	assert.Nil(t, Get("optional.value.but.not.found"))
}

func TestConfigValueRequiredDefault(t *testing.T) {
	err := setupTestConfig(schema.Config{
		{
			Key:      "required.value",
			Required: true,
			Default:  "default_fallback",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "default_fallback", Get("required.value"))
}
