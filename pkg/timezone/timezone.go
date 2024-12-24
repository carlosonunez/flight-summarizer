package timezone

import (
	"fmt"
	"time"
)

const ISO8601TimeFormat = "2006-01-02T15:04:05-07:00 MST"

// TimeZoneDatabase is a representation of a time zone database.
type TimeZoneDatabase interface {
	// LookupUTCOffsetByID retrieves the UTC offset given a three-digit identifier
	LookupUTCOffsetByID(ID string, start time.Time) (int64, error)
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

// MustParseISO8601Time parses a string in ISO8601 format; see also
// time.ISO8601TimeFormat. Panics if it fails.
func MustParseISO8601Time(tStr string) *time.Time {
	t, err := time.Parse(ISO8601TimeFormat, tStr)
	if err != nil {
		panic(fmt.Sprintf("time parsing during testing failed: %s", err))
	}
	return &t
}
