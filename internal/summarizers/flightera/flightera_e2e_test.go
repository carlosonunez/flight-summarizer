package flightera

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

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
		var wantJSON, gotJSON []byte
		s := FlighteraFlightSummarizer{}
		err := s.Init(&summarizer.FlightSummarizerOptions{FlightNumber: tc.FlightURLFragment})
		require.NoError(t, err, "test case: %s", tc.FlightURLFragment)
		want := tc.Want
		wantJSON, err = json.Marshal(&want)
		require.NoError(t, err)
		got, err := s.Summarize()
		require.NoError(t, err, "test case: %s", tc.FlightURLFragment)
		gotJSON, err = json.Marshal(&got)
		require.NoError(t, err)
		assert.Equal(t, string(wantJSON), string(gotJSON))
	}
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
