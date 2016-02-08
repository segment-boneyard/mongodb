package source_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/segmentio/go-source"
)

func TestSourceSetFalse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"id": 4,
			"result": false
		}`)
	}))
	defer ts.Close()

	source, err := source.New(&source.Config{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	err = source.Set("users", "bill", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})

	assert.Equal(t, errors.New("Source.Set failed with false"), err)
}

func TestSourceReportErrorFalse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"id": 4,
			"result": false
		}`)
	}))
	defer ts.Close()

	source, err := source.New(&source.Config{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	err = source.ReportError("some error", "users", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})

	assert.Equal(t, errors.New("Source.ReportError failed with false"), err)
}

func TestReportError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"id": 4,
			"result": true
		}`)
	}))
	defer ts.Close()

	source, err := source.New(&source.Config{
		URL: ts.URL,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = source.ReportError("users", "bill", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done!")

	// Output: Done!
}

func TestSet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"id": 4,
			"result": true
		}`)
	}))
	defer ts.Close()

	source, err := source.New(&source.Config{
		URL: ts.URL,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = source.Set("users", "bill", map[string]interface{}{
		"first_name": "Bill",
		"last_name":  "Lumbergh",
		"email":      "bill.lumbergh@initech.com",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Done!")

	// Output: Done!
}
