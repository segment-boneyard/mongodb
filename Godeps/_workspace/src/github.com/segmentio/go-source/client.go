package source

// Client wraps calls to the RPC service exposed by the source-runner in a Go API.
// Use `New` to create a client.
type Client interface {
	// Set an object with the given collection, id and properties.
	Set(collection string, id string, properties map[string]interface{}) error

	// Report an error, with an optional collection and properties.
	ReportError(message, collection string, properties map[string]interface{}) error
}

// Config wraps options that can be passed to the client.
type Config struct {
	URL string // URL of the source-runner RPC service.
}

// New creates a client instance with the given configuration.
func New(config *Config) (Client, error) {
	return newClient(config)
}
