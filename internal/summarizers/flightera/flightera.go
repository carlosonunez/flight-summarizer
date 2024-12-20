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
	summary := summarizer.FlightSummary{}
	originAirport, err := flighteraGetOriginAirport(b)
	if err != nil {
		return summary, error
	}
	destAirport, err := flighteraGetDestinationAirport(b)
	if err != nil {
		return summary, error
	}
	summary.FlightNumber = "WIP"
	return nil, errors.New("WIP")
}

// NewFlighteraFlightSummarizer creates a flight summarizer pulled from
// Flightera data.
func NewFlighteraFlightSummarizer(b browser.Browser) (*FlighteraFlightSummarizer, error) {
	return nil, errors.New("WIP")
}
