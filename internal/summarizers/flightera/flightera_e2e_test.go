package flightera

import (
	"testing"

	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlighteraLiveFlightsE2E(t *testing.T) {
	tcs, err := summarizer.LoadTestCases("tests.yaml")
	require.NoError(t, err)

	for _, tc := range tcs.LiveFlights {
		flightNumber := tc.Want.FlightNumber
		got, err := GenerateFlightSummary(flightNumber)
		require.NoError(t, err)
		assert.Equal(t, tc.Want, got)
	}
}
