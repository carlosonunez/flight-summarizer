package internal

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/carlosonunez/flight-summarizer/types"
	"golang.org/x/net/html"
)

type flightSideType int64

const (
	origin flightSideType = iota
	destination
)

type dateExtractType int64

const (
	monthAndDay dateExtractType = iota
	localTZ
	timeUTC
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

func flighteraGetOriginAirport(b types.Browser) (string, error) {
	return matchAirportIATA(b, origin)
}

func flighteraGetDestinationAirport(b types.Browser) (string, error) {
	return matchAirportIATA(b, destination)
}

func flighteraGetScheduledDeparture(b types.Browser, db types.TimeZoneDatabase, t flightSideType) (*time.Time, error) {
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
	monthAndDay, err := extractMonthDayFromRawTime(timeRaw)
	if err != nil {
		return nil, err
	}
	localTZ, err := extractLocalTZFromRawTime(timeRaw)
	if err != nil {
		return nil, err
	}
	timeUTC, err := extractTimeUTCFromRawTime(timeRaw)
	if err != nil {
		return nil, err
	}
	timeParsed, err := time.Parse("02 Jan 2006 15:04 UTC", fmt.Sprintf("%s %d %s",
		monthAndDay,
		year,
		timeUTC))
	if err != nil {
		return nil, err
	}
	offset, err := db.LookupUTCOffsetByID(localTZ, timeParsed)
	if err != nil {
		return nil, err
	}
	fixedTime := timeParsed.In(time.FixedZone(localTZ, int(offset)))
	fmt.Printf("parsed: %+v, fixed: %+v", timeParsed, fixedTime)
	return &fixedTime, nil
}

func extractMonthDayFromRawTime(text string) (string, error) {
	return extractFromRawTime(text, monthAndDay)
}

func extractLocalTZFromRawTime(text string) (string, error) {
	return extractFromRawTime(text, localTZ)
}

func extractTimeUTCFromRawTime(text string) (string, error) {
	return extractFromRawTime(text, timeUTC)
}

func extractFromRawTime(text string, t dateExtractType) (string, error) {
	var pattern string
	switch t {
	case monthAndDay:
		pattern = "[0-9]{1,2} [A-Za-z]{3}"
	case localTZ:
		pattern = ".*([A-Z]{3})$"
	case timeUTC:
		pattern = "(.*) UTC$"
	default:
		return "", fmt.Errorf("invalid date extract type: %d", t)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(text, "\n") {
		if re.Match([]byte(line)) {
			return strings.TrimSpace(re.FindAllString(line, -1)[0]), nil
		}
	}
	return "", fmt.Errorf("no time fragment found that matches expr '%s'", pattern)
}

func matchAirportIATA(b types.Browser, t flightSideType) (string, error) {
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

func getNodesOnPage(b types.Browser, query string) ([]*html.Node, error) {
	found, err := htmlquery.QueryAll(b.Document(), query)
	if err != nil {
		return found, nil
	}
	if len(found) == 0 {
		return found, fmt.Errorf("nothing match xpath: %s", query)
	}
	return found, nil
}

func getTextOnPage(b types.Browser, query string) ([]string, error) {
	return getTextOnPageRegexp(b, query, "")
}

func getTextOnPageRegexp(b types.Browser, query string, pattern string) ([]string, error) {
	var results []string
	found, err := getNodesOnPage(b, query)
	if err != nil {
		return results, err
	}
	re := regexp.MustCompile(pattern)
	results = re.FindAllString(htmlquery.InnerText(found[0]), -1)
	return results, nil
}
