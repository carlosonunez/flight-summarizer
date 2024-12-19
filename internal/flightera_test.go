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
	assert.Equal(t, iata, "CLT")
}

func TestDestinationAirport(t *testing.T) {
	b, err := NewMockBrowser("FAKE1")
	require.NoError(t, err)
	iata, err := flighteraGetDestinationAirport(b)
	require.NoError(t, err)
	assert.Equal(t, iata, "TUL")
}

func TestOriginScheduledDepartureTime(t *testing.T) {
	b, err := NewMockBrowser("FAKE1")
	require.NoError(t, err)
	tzdb, err := timezonedb.NewTimeZoneDBDotComDB(testhelpers.TIMEZONE_DB_FIXTURE_PATH)
	require.NoError(t, err)
	want := testhelpers.MustParseTime("2024-12-15T11:32:00-05:00")
	got, err := flighteraGetScheduledDeparture(b, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

// TODO: Finish this.
func TestOriginScheduledDeparture(t *testing.T) {
	t.Fail()
}
