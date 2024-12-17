package types

import "time"

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
	Scheduled *time.Time `json:"scheduled" yaml:"scheduled"`
	// Actual are actual (realistic) times.
	Actual *time.Time `json:"actual" yaml:"actual"`
}
