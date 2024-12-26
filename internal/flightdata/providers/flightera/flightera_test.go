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
var browserLiveOnTimeEarly, browserLiveLateEarly, browserScheduled *MockBrowser

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

func TestDestinationScheduledLandingTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-19T16:05:00-07:00 MST")
	got, err := GetDestinationScheduledLandingTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDestinationActualLandingTime(t *testing.T) {
	scheduled := timezone.MustParseISO8601Time("2024-12-19T16:05:00-07:00 MST")
	actual := timezone.MustParseISO8601Time("2024-12-19T15:52:00-07:00 MST")
	got, err := GetDestinationActualLandingTime(browserLiveLateEarly, tzdb)
	require.NoError(t, err)
	assert.NotEqual(t, scheduled, got)
	assert.Equal(t, actual, got)
}

func TestScheduledOriginScheduledDepartureTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-26T17:05:00-06:00 CST")
	got, err := GetOriginScheduledDepartureTime(browserScheduled, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestScheduledOriginActualDepartureTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-26T17:05:00-06:00 CST")
	got, err := GetOriginActualDepartureTime(browserScheduled, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestScheduledDestinationScheduledLandingTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-26T17:56:00-06:00 CST")
	got, err := GetDestinationScheduledLandingTime(browserScheduled, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestScheduledDestinationActualLandingTime(t *testing.T) {
	want := timezone.MustParseISO8601Time("2024-12-26T17:56:00-06:00 CST")
	got, err := GetDestinationActualLandingTime(browserScheduled, tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFlighteraFlightNumber(t *testing.T) {
	want := "AAL5005"
	got, err := GetFlightNumber(browserLiveLateEarly)
	require.NoError(t, err)
	assert.Equal(t, want, got)
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
	browserScheduled, err = NewMockBrowser(Scheduled)
	if err != nil {
		panic(err)
	}
	tzdb, err = timezonedb.NewTimeZoneDBDotComDB(&timezonedb.TimeZoneDBDotComDBOptions{
		CSVFile: "fixtures/timezonedb.csv",
	})
	if err != nil {
		panic(err)
	}
}
