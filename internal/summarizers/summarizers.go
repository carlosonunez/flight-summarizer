package summarizers

import (
	"fmt"
	"strings"

	"github.com/carlosonunez/flight-summarizer/internal/summarizers/flightera"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
)

var registeredSummarizers = map[string]summarizer.FlightSummarizer{
	"flightera": &flightera.FlighteraFlightSummarizer{},
	"test":      &summarizer.ExampleFlightSummarizer{},
}

var defaultSummarizer = "flightera"

// Lookup returns a registered summarizer.
func Lookup(key string) (summarizer.FlightSummarizer, error) {
	s, ok := registeredSummarizers[strings.ToLower(key)]
	if !ok {
		return s, fmt.Errorf("No summarizer found that matches key: %s", key)
	}
	return s, nil
}

// LookupDefault returns the default registered summarizer.
func LookupDefault() (summarizer.FlightSummarizer, error) {
	return Lookup(defaultSummarizer)
}
