package browser

import "golang.org/x/net/html"

// BrowserOpts let you provide custom options to a browser.  Provide
// `browser.NO_BROWSER_OPTIONS` when calling Visit() if the browser that you're
// using does not expose any options.
type BrowserOpts map[string]interface{}

var NO_BROWSER_OPTIONS = BrowserOpts{}

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
