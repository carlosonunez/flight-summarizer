package internal

import (
	"testing"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestinationAirport(t *testing.T) {
	b := MockBrowser{}
	err := b.Init("FAKE1")
	require.NoError(t, err)
	iata, err := flighteraGetDestinationAirport(&b)
	require.NoError(t, err)
	assert.Equal(t, iata, "TUL")
}

func TestOriginScheduledDepartureTime(t *testing.T) {
	b := MockBrowser{}
	err := b.Init("FAKE1")
	require.NoError(t, err)
	tzdb := MockTimeZoneDB{}
	err = tzdb.Init()
	require.NoError(t, err)
	want := testhelpers.MustParseTime("2024-12-15T11:32:00-05:00")
	got, err := flighteraGetScheduledDeparture(&b, &tzdb)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
