package summarizer

import (
	"fmt"
	"strings"
	"time"
)

type Time struct {
	*time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04 MST"))), nil
}

func (t *Time) MarshalYAML() ([]byte, error) {
	return []byte(t.Format("2006-01-02 15:04 MST")), nil
}

func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var s string
	if err = unmarshal(&s); err != nil {
		return err
	}
	formats := []string{
		"2006-01-02T15:04:05Z-07:00 MST",
		"2006-01-02T15:04:05Z-07:00",
		"2006-01-02T15:04:05-07:00",
	}
	var errors = make([]string, len(formats))
	for idx, format := range formats {
		if parsed, err := time.Parse(format, s); err == nil {
			*t = Time{&parsed}
			return nil
		} else {
			errors[idx] = fmt.Sprintf("  - %s", err)
		}
	}
	return fmt.Errorf("unable to parse time '%s':\n%+v", s, strings.Join(errors, "\n"))
}

// FlightSummarizer returns a FlightSummary from a byte array of data.
type FlightSummarizer interface {
	// Summarize does the thing!
	Summarize() (*FlightSummary, error)
}

// FlightSummary provides useful details about a flight
// without the junk.
type FlightSummary struct {
	// FlightNumber is the flight number with IATA airline identifiers.
	FlightNumber string `json:"flight_number" yaml:"flight_number"`

	// Origin provides information about the origin point for this route
	Origin *Point `json:"origin" yaml:"origin"`

	// Destination provides information about the origin point for this route
	Destination *Point `json:"destination" yaml:"destination"`
}

// Point provides a summary of a node in this route.
type Point struct {
	// Airport is the IATA identifier for the airport at this point.
	AirportIATA string `json:"airport" yaml:"airport"`

	// Times outlines the scheduled and actual times associated with this point.
	Times *FlightSummaryDateTimes `json:"times" yaml:"times"`
}

// FlightSummaryDateTimes provides scheduled and actual datetimes for this
// flight.
type FlightSummaryDateTimes struct {
	// Scheduled are scheduled (unrealistic) times.
	Scheduled *Time `json:"scheduled" yaml:"scheduled"`
	// Actual are actual (realistic) times.
	Actual *Time `json:"actual" yaml:"actual"`
}

// NewEmptyFlightSummary creates an empty flight summary (so that the summarizer
// doesn't have to type out this ugly struct initializer themselves)
func NewEmptyFlightSummary() *FlightSummary {
	return &FlightSummary{
		Origin: &Point{
			Times: &FlightSummaryDateTimes{
				Scheduled: &Time{},
				Actual:    &Time{},
			},
		},
		Destination: &Point{
			Times: &FlightSummaryDateTimes{
				Scheduled: &Time{},
				Actual:    &Time{},
			},
		},
	}
}
