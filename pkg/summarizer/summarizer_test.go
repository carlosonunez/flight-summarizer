package summarizer

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummaryToJSON(t *testing.T) {
	wantRaw := `
{
	"flight_number": "FAKE1",
	"origin": {
		"airport": "FOO",
		"city": "Foo City",
		"times": {
				"scheduled": "2024-12-22 17:22 PST",
				"actual": "2024-12-22 17:00 PST"
			}
	},
	"destination": {
		"airport": "BAR",
		"city": "Bar City",
		"times": {
			"scheduled": "2024-12-22 18:22 PST",
			"actual": "2024-12-22 18:00 PST"
		}
	}
}`
	wantObj := FlightSummary{
		FlightNumber: "FAKE1",
		Origin: &Point{
			AirportIATA: "FOO",
			City:        "Foo City",
			Times: &FlightSummaryDateTimes{
				Scheduled: &Time{mustParseTime("2024-12-22 17:22 PST")},
				Actual:    &Time{mustParseTime("2024-12-22 17:00 PST")},
			},
		},
		Destination: &Point{
			AirportIATA: "BAR",
			City:        "Bar City",
			Times: &FlightSummaryDateTimes{
				Scheduled: &Time{mustParseTime("2024-12-22 18:22 PST")},
				Actual:    &Time{mustParseTime("2024-12-22 18:00 PST")},
			},
		},
	}
	var want bytes.Buffer
	require.NoError(t, json.Indent(&want, []byte(wantRaw), "", "	"))
	got, err := json.MarshalIndent(&wantObj, "", "	")
	require.NoError(t, err)
	assert.Equal(t, string(want.Bytes()), string(got))
}

func mustParseTime(t string) *time.Time {
	parsed, err := time.Parse("2006-01-02 15:04 MST", t)
	if err != nil {
		panic(err)
	}
	return &parsed
}
