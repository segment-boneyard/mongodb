package mock_test

import (
	"errors"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/segmentio/go-source"
	"github.com/segmentio/go-source/mock"
)

// assert interface compliance.
var _ source.Client = (*mock.Source)(nil)

func TestMockSet(t *testing.T) {
	source := mock.New()

	err := source.Set("users", "bill", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, <-source.SetCalls, mock.Set{
		Collection: "users",
		ID:         "bill",
		Properties: map[string]interface{}{
			"first_name": "Bill",
			"last_name":  "Lumbergh",
			"email":      "bill.lumbergh@initech.com",
		},
	})
}

func TestMockReportError(t *testing.T) {
	source := mock.New()

	err := source.ReportError("panic", "users", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, <-source.ReportErrorCalls, mock.ReportError{
		Message:    "panic",
		Collection: "users",
		Properties: map[string]interface{}{
			"first_name": "Bill",
			"last_name":  "Lumbergh",
			"email":      "bill.lumbergh@initech.com",
		},
	})
}

func TestSetErrors(t *testing.T) {
	source := mock.New()

	{
		err := source.Set("", "bill", map[string]interface{}{
			"first_name": "Bill",
		})
		assert.Equal(t, errors.New("Parameter `collection` must not be empty."), err)
	}
	{
		err := source.Set("users", "", map[string]interface{}{
			"first_name": "Bill",
		})
		assert.Equal(t, errors.New("Parameter `id` must not be empty."), err)
	}
	{
		err := source.Set("users", "billd", map[string]interface{}{})
		assert.Equal(t, errors.New("Parameter `properties` must not be empty."), err)
	}
}

func TestReportErrorErrors(t *testing.T) {
	source := mock.New()

	err := source.ReportError("", "bill", map[string]interface{}{
		"first_name": "Bill",
	})
	assert.Equal(t, errors.New("Parameter `message` must not be empty."), err)
}
