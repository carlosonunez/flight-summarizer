package flightera

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type testCaseList struct {
	Tests []*testCase `yaml:"tests"`
}

type testCase struct {
	FlightURLFragment string                    `yaml:"stub"`
	Want              *summarizer.FlightSummary `yaml:"want"`
}

var tcs *testCaseList

func TestFlighteraLiveFlightsE2E(t *testing.T) {
	for _, tc := range tcs.Tests {
		s := FlighteraFlightSummarizer{}
		err := s.Init(&summarizer.FlightSummarizerOptions{FlightNumber: tc.FlightURLFragment})
		require.NoError(t, err, "test case: %s", tc.FlightURLFragment)
		want := tc.Want
		got, err := s.Summarize()
		require.NoError(t, err, "test case: %s", tc.FlightURLFragment)
		assert.Equal(t, want.Origin.AirportIATA, got.Origin.AirportIATA, "[origin] airports match: %s", tc.FlightURLFragment)
		assert.Equal(t, want.Destination.AirportIATA, got.Destination.AirportIATA, "[dest] airports match: %s", tc.FlightURLFragment)
		testTimes(t, want.Origin.Times.Scheduled, got.Origin.Times.Scheduled, "[origin] scheduled times match: %s", tc.FlightURLFragment)
		testTimes(t, want.Origin.Times.Actual, got.Origin.Times.Actual, "[origin] actual times match: %s", tc.FlightURLFragment)
		testTimes(t, want.Destination.Times.Scheduled, got.Destination.Times.Scheduled, "[dest] scheduled times match: %s", tc.FlightURLFragment)
		testTimes(t, want.Destination.Times.Actual, got.Destination.Times.Actual, "[dest] actual times match: %s", tc.FlightURLFragment)
	}
}

func testTimes(t *testing.T, want *summarizer.Time, got *summarizer.Time, format string, parts ...any) {
	require.NotNil(t, got, fmt.Sprintf(format, parts...))
	// time.Time.Equal(...) will be false for times that don't have the same
	// location data in them, which ours do not.
	// We're going to do string-based equality to work around this.
	wantStr := want.Format(time.RFC3339)
	gotStr := got.Format(time.RFC3339)
	assert.Equal(t, wantStr, gotStr, fmt.Sprintf(format, parts...))
}

func init() {
	f, err := os.ReadFile("e2e/tests.yaml")
	if err != nil {
		panic(fmt.Sprintf("couldn't find flightera e2e test definitions: %s", err))
	}
	if err := yaml.Unmarshal(f, &tcs); err != nil {
		panic(fmt.Sprintf("couldn't load flightera test cases: %s", err))
	}

}
