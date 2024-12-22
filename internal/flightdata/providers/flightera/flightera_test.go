package flightera

import (
	"net/url"
	"testing"

	"github.com/carlosonunez/flight-summarizer/internal/timezone/providers/timezonedb"
	"github.com/carlosonunez/flight-summarizer/pkg/timezone"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tzdb *timezonedb.TimeZoneDBDotComDB
var browserLiveOnTimeEarly, browserLiveLateEarly *MockBrowser

func TestOriginAirport(t *testing.T) {
	iata, err := GetOriginAirport(browserLiveOnTimeEarly)
	require.NoError(t, err)
	assert.Equal(t, "DFW", iata)
}

func TestDestinationAirport(t *testing.T) {
	iata, err := GetDestinationAirport(browserLiveOnTimeEarly)
	require.NoError(t, err)
	assert.Equal(t, "LGA", iata)
}

func TestOriginScheduledDepartureTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-19T14:52:00-06:00 CST")
	got, err := GetOriginScheduledDepartureTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestOriginActualDepartureTime(t *testing.T) {
	actual := timezone.MustParseISO8601Time("2024-12-19T15:10:00-06:00 CST")
	scheduled := timezone.MustParseISO8601Time("2024-12-19T14:52:00-06:00 CST")
	got, err := GetOriginActualDepartureTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.NotEqual(t, scheduled, got)
	assert.Equal(t, actual, got)
}

func TestDestinationScheduledDepartureTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-19T14:52:00-06:00 CST")
	got, err := GetDestinationScheduledDepartureTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDestinationActualDepartureTime(t *testing.T) {
	actual := timezone.MustParseISO8601Time("2024-12-19T15:10:00-06:00 CST")
	scheduled := timezone.MustParseISO8601Time("2024-12-19T14:52:00-06:00 CST")
	got, err := GetDestinationActualDepartureTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.NotEqual(t, scheduled, got)
	assert.Equal(t, actual, got)
}

func TestFlighteraURLJustFlightNumber(t *testing.T) {
	want, err := url.Parse("https://flightera.net/en/flight/FAKE1")
	require.NoError(t, err)
	got, err := flighteraURL("FAKE1")
	require.NoError(t, err)
	assert.Equal(t, want, got)
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
	tzdb, err = timezonedb.NewMockTimeZoneDBDotComDB()
	if err != nil {
		panic(err)
	}
}
