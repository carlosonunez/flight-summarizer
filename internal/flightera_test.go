package internal

import (
	"testing"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/carlosonunez/flight-summarizer/timezonedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tzdb *timezonedb.TimeZoneDBDotComDB
var browserLiveOnTimeEarly, browserLiveLateEarly *MockBrowser

func TestOriginAirport(t *testing.T) {
	iata, err := flighteraGetOriginAirport(browserLiveOnTimeEarly)
	require.NoError(t, err)
	assert.Equal(t, "DFW", iata)
}

func TestDestinationAirport(t *testing.T) {
	iata, err := flighteraGetDestinationAirport(browserLiveOnTimeEarly)
	require.NoError(t, err)
	assert.Equal(t, "LGA", iata)
}

func TestOriginScheduledDepartureTime(t *testing.T) {
	want := testhelpers.MustParseTime("2024-12-19T14:52:00-06:00 CST")
	got, err := flighteraGetScheduledDeparture(browserLiveLateEarly, tzdb, origin)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestOriginActualDepartureTime(t *testing.T) {
	actual := testhelpers.MustParseTime("2024-12-19T15:10:00-06:00 CST")
	scheduled := testhelpers.MustParseTime("2024-12-19T14:52:00-06:00 CST")
	got, err := flighteraGetActualDeparture(browserLiveLateEarly, tzdb, origin)
	require.NoError(t, err)
	assert.NotEqual(t, scheduled, got)
	assert.Equal(t, actual, got)
}

func TestDestinationScheduledDepartureTime(t *testing.T) {
	tzdb, err := timezonedb.NewTimeZoneDBDotComDB(testhelpers.TIMEZONE_DB_FIXTURE_PATH)
	require.NoError(t, err)
	want := testhelpers.MustParseTime("2024-12-19T14:52:00-06:00 CST")
	got, err := flighteraGetScheduledDeparture(browserLiveLateEarly, tzdb, origin)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDestinationActualDepartureTime(t *testing.T) {
	actual := testhelpers.MustParseTime("2024-12-19T15:10:00-06:00 CST")
	scheduled := testhelpers.MustParseTime("2024-12-19T14:52:00-06:00 CST")
	got, err := flighteraGetActualDeparture(browserLiveLateEarly, tzdb, origin)
	require.NoError(t, err)
	assert.NotEqual(t, scheduled, got)
	assert.Equal(t, actual, got)
}

func init() {
	var err error
	browserLiveOnTimeEarly, err = NewMockBrowser(LiveOnTimeDepartureEarlyArrival)
	if err != nil {
		panic(err)
	}
	browserLiveLateEarly, err = NewMockBrowser(LiveLateDepartureEarlyArrival)
	if err != nil {
		panic(err)
	}
	tzdb, err = timezonedb.NewTimeZoneDBDotComDB(testhelpers.TIMEZONE_DB_FIXTURE_PATH)
	if err != nil {
		panic(err)
	}
}
