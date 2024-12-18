package internal

import "github.com/carlosonunez/flight-summarizer/types"

// FlighteraFlightSummarizer summarizes flight summaries from Flightera.
type FlighteraFlightSummarizer struct{}

func (s *FlighteraFlightSummarizer) Summarize(data []byte) (*types.FlightSummary, error) {
	return nil, nil
}

// GenerateFlightSummary generates a flight summary.
func GenerateFlightSummary(flightNumber string) (*types.FlightSummary, error) {
	return nil, nil
}
