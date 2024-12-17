package internal

import (
	"testing"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLiveFlightsE2E(t *testing.T) {
	tcs, err := testhelpers.LoadTestCases()
	require.NoError(t, err)

	for _, tc := range tcs.LiveFlights {
		flightNumber := tc.Want.FlightNumber
		got, err := GenerateFlightSummary(flightNumber)
		require.NoError(t, err)
		assert.Equal(t, tc.Want, got)
	}
}
