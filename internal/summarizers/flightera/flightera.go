package flightera

import (
	"errors"

	"github.com/carlosonunez/flight-summarizer/pkg/browser"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
)

// FlighteraFlightSummarizer summarizes flights from Flightera, a lightweight
// flight data provider.
type FlighteraFlightSummarizer struct{}

// Summarize does the summarization.
func (s *FlighteraFlightSummarizer) Summarize(b browser.Browser) (*summarizer.FlightSummary, error) {
	return nil, errors.New("WIP")
}

// NewFlighteraFlightSummarizer creates a flight summarizer pulled from
// Flightera data.
func NewFlighteraFlightSummarizer(b browser.Browser) (*FlighteraFlightSummarizer, error) {
	return nil, errors.New("WIP")
}
