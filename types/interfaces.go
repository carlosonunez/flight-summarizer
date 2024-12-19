package types

import (
	"time"

	"golang.org/x/net/html"
)

// BrowserOpts let you provide custom options to a browser.
type BrowserOpts map[string]interface{}

// Browser is an interface that's used to retrieve
// online flight data. This makes it easier to retrieve
// data from text-based and "full" headless browsers depending on
// the service.
type Browser interface {
	// Init initializes the Browser.
	Init(flightNumber string) error

	// Visit visits a URL managed by the Browser-compliant object, optionally with
	// some options to change its behavior.
	Visit(opts BrowserOpts) error

	// Document returns the contents of the page for processing in a FlightSummarizer.
	Document() *html.Node
}

// FlightSummarizer returns a FlightSummary from a byte array of data.
type FlightSummarizer interface {
	// Summarize does the thing!
	Summarize(data []byte) (*FlightSummary, error)
}

// TimeZoneDatabase is a representation of a time zone database.
type TimeZoneDatabase interface {
	// LookupUTCOffsetByID retrieves the UTC offset given a three-digit identifier
	LookupUTCOffsetByID(ID string, start time.Time) (int64, error)
}
