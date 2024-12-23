package flightera

import (
	"github.com/carlosonunez/flight-summarizer/internal/flightdata/providers/flightera"
	"github.com/carlosonunez/flight-summarizer/internal/timezone/providers/timezonedb"
	"github.com/carlosonunez/flight-summarizer/pkg/browser"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
	"github.com/carlosonunez/flight-summarizer/pkg/timezone"
)

// FlighteraFlightSummarizer summarizes flights from Flightera, a lightweight
// flight data provider.
type FlighteraFlightSummarizer struct {
	browser browser.Browser
	tzdb    timezone.TimeZoneDatabase
	summary *summarizer.FlightSummary
}

// Summarize does the summarization.
func (s *FlighteraFlightSummarizer) Summarize() (*summarizer.FlightSummary, error) {
	b := s.browser
	ops := []func(*FlighteraFlightSummarizer, browser.Browser) error{
		retrieveFlightNumber,
		summarizeAirports,
		summarizeDepartureTimes,
		summarizeLandingTimes,
	}
	for _, op := range ops {
		if err := op(s, b); err != nil {
			return s.summary, err
		}
	}
	return s.summary, nil
}

func retrieveFlightNumber(s *FlighteraFlightSummarizer, b browser.Browser) (err error) {
	if s.summary.FlightNumber, err = flightera.GetFlightNumber(b); err != nil {
		return err
	}
	return nil
}

func summarizeAirports(s *FlighteraFlightSummarizer, b browser.Browser) (err error) {
	if s.summary.Origin.AirportIATA, err = flightera.GetOriginAirport(b); err != nil {
		return err
	}
	if s.summary.Destination.AirportIATA, err = flightera.GetDestinationAirport(b); err != nil {
		return err
	}
	return nil
}

func summarizeDepartureTimes(s *FlighteraFlightSummarizer, b browser.Browser) (err error) {
	if s.summary.Origin.Times.Scheduled, err = flightera.GetOriginScheduledDepartureTime(b, s.tzdb); err != nil {
		return err
	}
	if s.summary.Origin.Times.Actual, err = flightera.GetOriginActualDepartureTime(b, s.tzdb); err != nil {
		return err
	}
	return nil
}

func summarizeLandingTimes(s *FlighteraFlightSummarizer, b browser.Browser) (err error) {
	if s.summary.Destination.Times.Scheduled, err = flightera.GetDestinationScheduledLandingTime(b, s.tzdb); err != nil {
		return err
	}
	if s.summary.Destination.Times.Actual, err = flightera.GetDestinationActualLandingTime(b, s.tzdb); err != nil {
		return err
	}
	return nil
}

// NewFlighteraFlightSummarizer creates a flight summarizer pulled from
// Flightera data retrieved via a text-based browser.
func NewFlighteraFlightSummarizer(flightNumber string) (*FlighteraFlightSummarizer, error) {
	var err error
	s := FlighteraFlightSummarizer{summary: summarizer.NewEmptyFlightSummary()}
	if s.browser, err = flightera.NewFlighteraTextBrowser(flightNumber); err != nil {
		return nil, err
	}
	if err = s.browser.Visit(browser.NO_BROWSER_OPTIONS); err != nil {
		return nil, err
	}
	if s.tzdb, err = timezonedb.NewTimeZoneDBDotComDB(timezonedb.DEFAULT_OPTIONS); err != nil {
		return nil, err
	}
	return &s, nil
}
