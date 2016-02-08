package mock

import "errors"

// Set is Go struct representation of a `Source.Set` call.
type Set struct {
	Collection string
	ID         string
	Properties map[string]interface{}
}

// ReportError is Go struct representation of a `Source.ReportError` call.
type ReportError struct {
	Message    string
	Collection string
	Properties map[string]interface{}
}

// Source is a mock source client that can be used to record set calls for
// inspection in tests.
type Source struct {
	SetCalls         chan Set
	ReportErrorCalls chan ReportError
}

func New() *Source {
	return &Source{
		SetCalls:         make(chan Set, 10),
		ReportErrorCalls: make(chan ReportError, 10),
	}
}

func (s *Source) Set(collection, id string, properties map[string]interface{}) error {
	if collection == "" {
		return errors.New("Parameter `collection` must not be empty.")
	}
	if id == "" {
		return errors.New("Parameter `id` must not be empty.")
	}
	if len(properties) == 0 {
		return errors.New("Parameter `properties` must not be empty.")
	}
	s.SetCalls <- Set{collection, id, properties}
	return nil
}

func (s *Source) ReportError(message, collection string, properties map[string]interface{}) error {
	if message == "" {
		return errors.New("Parameter `message` must not be empty.")
	}
	s.ReportErrorCalls <- ReportError{message, collection, properties}
	return nil
}
