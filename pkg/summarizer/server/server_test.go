package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrievesFlightSummary(t *testing.T) {
	var want bytes.Buffer
	wantRaw := `{
	"flight_number": "FAKE1",
	"origin": {
		"airport": "FOO",
		"times": {
				"scheduled": "2024-12-23 11:07 PST",
				"actual": "2024-12-23 11:00 PST"
			}
	},
	"destination": {
		"airport": "BAR",
		"times": {
				"scheduled": "2024-12-23 11:07 PST",
				"actual": "2024-12-23 11:00 PST"
		}
	}
}`
	require.NoError(t, json.Indent(&want, []byte(wantRaw), "", "	"))
	req := httptest.NewRequest(http.MethodGet, "/summarize?flightNumber=FOOBAR1&summarizer=test", nil)
	w := httptest.NewRecorder()
	summarizeHandler(w, req)
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	defer res.Body.Close()
	gotRaw, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	var got bytes.Buffer
	assert.NoError(t, json.Indent(&got, gotRaw, "", "	"))
	assert.Equal(t, want.String(), got.String())
}
