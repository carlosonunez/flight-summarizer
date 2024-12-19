package internal

import (
	"testing"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/carlosonunez/flight-summarizer/timezonedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOriginAirport(t *testing.T) {
	b, err := NewMockBrowser("FAKE1")
	require.NoError(t, err)
	iata, err := flighteraGetOriginAirport(b)
	require.NoError(t, err)
	assert.Equal(t, "DFW", iata)
}

func TestDestinationAirport(t *testing.T) {
	b, err := NewMockBrowser("FAKE1")
	require.NoError(t, err)
	iata, err := flighteraGetDestinationAirport(b)
	require.NoError(t, err)
	assert.Equal(t, "LGA", iata)
}

func TestOriginScheduledDepartureTime(t *testing.T) {
	b, err := NewMockBrowser("FAKE1")
	require.NoError(t, err)
	tzdb, err := timezonedb.NewTimeZoneDBDotComDB(testhelpers.TIMEZONE_DB_FIXTURE_PATH)
	require.NoError(t, err)
	want := testhelpers.MustParseTime("2024-12-19T14:53:00-06:00 CST")
	got, err := flighteraGetScheduledDeparture(b, tzdb, origin)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
