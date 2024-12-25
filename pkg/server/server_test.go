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

func TestReadyEndpoint(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	require.NoError(t, err)
	pingHandler(w, req)
	testStatusCodeOK(t, w.Result())
	testBodyEqual(t, `{"status":"ok"}`, w.Result())
}

func testStatusCodeOK(t *testing.T, r *http.Response) {
	assert.Equal(t, http.StatusOK, r.StatusCode)
}

func testBodyEqual(t *testing.T, wantRaw string, r *http.Response) {
	var want bytes.Buffer
	var got bytes.Buffer
	require.NoError(t, json.Indent(&want, []byte(wantRaw), "", "	"))
	defer r.Body.Close()
	gotRaw, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	assert.NoError(t, json.Indent(&got, gotRaw, "", "	"))
	assert.Equal(t, want.String(), got.String())
}

func TestRetrievesFlightSummary(t *testing.T) {
	wantRaw := `{
  "status": "ok",
	"summary": {
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
	}
}`
	req := httptest.NewRequest(http.MethodGet, "/summarize?flightNumber=FOOBAR1&summarizer=test", nil)
	w := httptest.NewRecorder()
	summarizeHandler(w, req)
	res := w.Result()
	testStatusCodeOK(t, res)
	testBodyEqual(t, wantRaw, res)
}
