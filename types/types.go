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

// TimeZoneDBInfo is a time zone entry in a time zone database from
// timezonedb.com
type TimeZoneDBInfo struct {
	// IsDaylightSavingsTime indicates whether DST applies for this entry.
	IsDaylightSavingsTime int `csv:"dst"`
	// ZoneName is the friendly name of the zone, like 'America/Houston'
	// TimeStart is the time at which the Offset corresponding to this entry
	// applies. This database has multiple time samples to properly account for
	// daylight savings time.
	TimeStart int32 `csv:"time_start"`
	// UTCOffsetSeconds is the GMT/UTC offset, in seconds.
	UTCOffsetSeconds int32  `csv:"gmt_offset"`
	ZoneName         string `csv:"zone_name"`
	// CountryCode is a two-char ID of the country the zone belongs to, like 'US'
	CountryCode string `csv:"country_code"`
	// Abbreviation is the timezone abbreviation for the zone at a given time,
	// like 'CST'
	Abbreviation string `csv:"abbreviation"`
}
