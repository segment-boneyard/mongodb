package source

import (
	"errors"

	jsonrpc "github.com/gohttp/jsonrpc-client"
)

type client struct {
	jsonrpcClient *jsonrpc.Client
}

func newClient(config *Config) (*client, error) {
	jsonrpcClient := jsonrpc.NewClient(config.URL)

	return &client{
		jsonrpcClient: jsonrpcClient,
	}, nil
}

// todo: reuse from source-runner pkg?
type setRequest struct {
	ID         string
	Properties map[string]interface{}
	Collection string
}

func (c *client) Set(collection string, id string, properties map[string]interface{}) error {
	req := &setRequest{
		ID:         id,
		Properties: properties,
		Collection: collection,
	}

	var result bool
	err := c.jsonrpcClient.Call("Source.Set", req, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Source.Set failed with false")
	}

	return nil
}

type reportErrorRequest struct {
	Message    string                 `json:"message"`
	Collection string                 `json:"collection,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func (c *client) ReportError(message string, collection string, properties map[string]interface{}) error {
	req := &reportErrorRequest{
		Message:    message,
		Properties: properties,
		Collection: collection,
	}

	var result bool
	err := c.jsonrpcClient.Call("Source.ReportError", req, &result)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("Source.ReportError failed with false")
	}

	return nil
}
