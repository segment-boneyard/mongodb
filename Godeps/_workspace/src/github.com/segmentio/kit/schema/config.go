package schema

// Schema is the struct that defines your configuration
// settings, keys, required and optional values.
type Config []ConfigValue

// ConfigValue is the struct that contains information about a desired key.
//
// A ConfigValue is usually used inside a Schema, where a user can define
// before the application runs, which keys we should look for.
type ConfigValue struct {
	Key         string
	Required    bool
	Default     interface{}
	Description string
}
