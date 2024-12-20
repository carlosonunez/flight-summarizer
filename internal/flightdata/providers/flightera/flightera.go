package flightera

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/carlosonunez/flight-summarizer/pkg/browser"
	"github.com/carlosonunez/flight-summarizer/pkg/timezone"
	"golang.org/x/net/html"
)

type flightSideType int64
type datePartExtractType int64
type dateExtractType int64

const (
	origin flightSideType = iota
	destination
	monthAndDay datePartExtractType = iota
	localTZ
	localTime
	scheduled dateExtractType = iota
	actual
)

func (t flightSideType) String() string {
	switch t {
	case origin:
		return "origin"
	case destination:
		return "destination"
	default:
		return "????"
	}
}

func flighteraGetOriginAirport(b browser.Browser) (string, error) {
	return matchAirportIATA(b, origin)
}

func flighteraGetDestinationAirport(b browser.Browser) (string, error) {
	return matchAirportIATA(b, destination)
}

func flighteraGetScheduledDeparture(b browser.Browser, db timezone.TimeZoneDatabase, t flightSideType) (*time.Time, error) {
	return flighteraGetTime(b, db, t, scheduled)
}

func flighteraGetActualDeparture(b browser.Browser, db timezone.TimeZoneDatabase, t flightSideType) (*time.Time, error) {
	return flighteraGetTime(b, db, t, actual)
}

func flighteraGetTime(b browser.Browser, db timezone.TimeZoneDatabase, t flightSideType, dt dateExtractType) (*time.Time, error) {
	var idx int
	switch t {
	case origin:
		idx = 0
	case destination:
		idx = 1
	default:
		return nil, fmt.Errorf("no matching flight side type: %d", t)
	}
	nodes, err := getNodesOnPage(b, "//span[text()[contains(.,'UTC')]]")
	if err != nil {
		return nil, err
	}
	// Unfortunately, Flightera doesn't expose the year the flight occurred on the
	// page. We'll assume that it happened this year.
	// TODO: Use the year supplied by query param when Flight Summarizer supports
	// dates.
	year := time.Now().Year()
	timeRaw := htmlquery.InnerText(nodes[idx].Parent)
	monthAndDay, err := extractMonthDayFromRawTime(timeRaw, dt)
	if err != nil {
		return nil, err
	}
	localTZ, err := extractLocalTZFromRawTime(timeRaw, dt)
	if err != nil {
		return nil, err
	}
	localTime, err := extractLocalTimeFromRawTime(timeRaw, dt)
	if err != nil {
		return nil, err
	}
	timeParsed, err := time.Parse("02 Jan 2006 15:04 UTC", fmt.Sprintf("%s %d %s UTC",
		monthAndDay,
		year,
		localTime))
	if err != nil {
		return nil, err
	}
	// IANA's timezone identifiers (which is what Go uses for time.Location) are
	// incomplete. Given that we are fetching local time from these entries and
	// are assuming UTC, we'll need to manually offset the time to make sure that
	// it's correct.
	offset, err := db.LookupUTCOffsetByID(localTZ, timeParsed)
	if err != nil {
		return nil, err
	}
	fixedTime := timeParsed.Add(-time.Duration(offset) * time.Second).In(time.FixedZone(localTZ, int(offset)))
	return &fixedTime, nil
}

func extractMonthDayFromRawTime(text string, dt dateExtractType) (string, error) {
	return extractFromRawTime(text, monthAndDay, dt)
}

func extractLocalTZFromRawTime(text string, dt dateExtractType) (string, error) {
	return extractFromRawTime(text, localTZ, dt)
}

func extractLocalTimeFromRawTime(text string, dt dateExtractType) (string, error) {
	return extractFromRawTime(text, localTime, dt)
}

func stripJunkFromRawTime(text string) []string {
	re := regexp.MustCompile("([0-9]{1,2} [A-Za-z]+ [0-9]{2}:[0-9]{2}|[A-Z]{3}|[0-9]{2}:[0-9]{2} [A-Z]{3})")
	return re.FindAllString(text, -1)
}

func extractFromRawTime(text string, t datePartExtractType, dt dateExtractType) (string, error) {
	var pattern string
	var wantIndex int
	var scheduleDriftDetected, considerScheduleDrifts, seen bool
	timeParts := stripJunkFromRawTime(text)
	if len(timeParts) > 3 {
		scheduleDriftDetected = true
	}
	switch t {
	case monthAndDay:
		pattern = "[0-9]{1,2} [A-Za-z]{3}"
		if dt == scheduled {
			considerScheduleDrifts = true
		}
	case localTZ:
		pattern = "^([A-Z]{3})$"
	case localTime:
		wantIndex = 1
		pattern = "[0-9]{1,2} [A-Za-z]{3} ([0-9]{2}:[0-9]{2})"
		if dt == scheduled {
			considerScheduleDrifts = true
		}
	default:
		return "", fmt.Errorf("invalid date extract type: %d", t)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	for _, line := range timeParts {
		if re.Match([]byte(line)) {
			found := strings.TrimSpace(re.FindAllStringSubmatch(line, -1)[0][wantIndex])
			if !considerScheduleDrifts || (considerScheduleDrifts && !scheduleDriftDetected) {
				return found, nil
			}
			if scheduleDriftDetected && seen {
				return found, nil
			}
			seen = true
		}
	}
	return "", fmt.Errorf("no time fragment found that matches expr '%s'", pattern)
}

func matchAirportIATA(b browser.Browser, t flightSideType) (string, error) {
	var result string
	var query string
	switch t {
	case origin:
		query = "//div[contains(@class, \"items-start\")]"
	case destination:
		query = "//div[contains(@class, \"items-end\")]"
	default:
		return "", fmt.Errorf("no matching flight side type: %d", t)
	}
	found, err := getTextOnPageRegexp(b, query, "[A-Z]{3}")
	if err != nil {
		return result, fmt.Errorf("couldn't find %s airport: %s", t.String(), err)
	}
	return found[0], nil
}

func getNodesOnPage(b browser.Browser, query string) ([]*html.Node, error) {
	found, err := htmlquery.QueryAll(b.Document(), query)
	if err != nil {
		return found, nil
	}
	if len(found) == 0 {
		return found, fmt.Errorf("nothing match xpath: %s", query)
	}
	return found, nil
}

func getTextOnPage(b browser.Browser, query string) ([]string, error) {
	return getTextOnPageRegexp(b, query, "")
}

func getTextOnPageRegexp(b browser.Browser, query string, pattern string) ([]string, error) {
	var results []string
	found, err := getNodesOnPage(b, query)
	if err != nil {
		return results, err
	}
	re := regexp.MustCompile(pattern)
	results = re.FindAllString(htmlquery.InnerText(found[0]), -1)
	return results, nil
}
