package internal

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/carlosonunez/flight-summarizer/types"
)

type airportType int64

const (
	originAirport airportType = iota
	destinationAirport
)

func (t airportType) String() string {
	switch t {
	case originAirport:
		return "origin"
	case destinationAirport:
		return "destination"
	default:
		return "????"
	}
}

func flighteraGetOriginAirport(b types.Browser) (string, error) {
	return matchAirportIATA(b, originAirport)
}

func flighteraGetDestinationAirport(b types.Browser) (string, error) {
	return matchAirportIATA(b, destinationAirport)
}

func flighteraGetScheduledDeparture(b types.Browser, db types.TimeZoneDatabase) (*time.Time, error) {
	return nil, errors.New("WIP")
}

func matchAirportIATA(b types.Browser, t airportType) (string, error) {
	var result string
	var query string
	switch t {
	case originAirport:
		query = "//div[contains(@class, \"items-start\")]"
	case destinationAirport:
		query = "//div[contains(@class, \"items-end\")]"
	default:
		return "", fmt.Errorf("no matching airport type: %d", t)
	}
	found, err := matchFirstOnPageRegexp(b, query, "[A-Z]{3}")
	if err != nil {
		return result, fmt.Errorf("couldn't find %s airport: %s", t.String(), err)
	}
	return found[0], nil
}

func matchFirstOnPage(b types.Browser, query string) ([]string, error) {
	return matchFirstOnPageRegexp(b, query, "")
}

func matchFirstOnPageRegexp(b types.Browser, query string, pattern string) ([]string, error) {
	var results []string
	found, err := htmlquery.QueryAll(b.Document(), query)
	if err != nil {
		return results, nil
	}
	if len(found) == 0 {
		return results, fmt.Errorf("nothing match xpath: %s", query)
	}
	re := regexp.MustCompile(pattern)
	results = re.FindAllString(htmlquery.InnerText(found[0]), -1)
	return results, nil
}
