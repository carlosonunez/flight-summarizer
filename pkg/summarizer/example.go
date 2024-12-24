package summarizer

import "time"

// ExampleFlightSummarizer is a fake  Useful for creating new
// summarizers.
type ExampleFlightSummarizer struct {
	summary *FlightSummary
}

// Init initializes this fake
func (s *ExampleFlightSummarizer) Init(opts *FlightSummarizerOptions) error {
	s.summary = NewEmptyFlightSummary()
	return nil
}

// Summarize performs a fake summarization.
func (s *ExampleFlightSummarizer) Summarize() (*FlightSummary, error) {
	s.summary.FlightNumber = "FAKE1"
	s.summary.Origin.AirportIATA = "FOO"
	s.summary.Origin.Times.Scheduled.Time = mustParseRFC3339Time("2024-12-23T11:07:00Z-08:00 PST")
	s.summary.Origin.Times.Actual.Time = mustParseRFC3339Time("2024-12-23T11:00:00Z-08:00 PST")
	s.summary.Destination.AirportIATA = "BAR"
	s.summary.Destination.Times.Scheduled.Time = mustParseRFC3339Time("2024-12-23T11:07:00Z-08:00 PST")
	s.summary.Destination.Times.Actual.Time = mustParseRFC3339Time("2024-12-23T11:00:00Z-08:00 PST")
	return s.summary, nil
}

func mustParseRFC3339Time(t string) *time.Time {
	parsed, err := time.Parse("2006-01-02T15:04:05Z-07:00 MST", t)
	if err != nil {
		panic(err)
	}
	return &parsed
}
