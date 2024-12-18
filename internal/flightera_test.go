package internal

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/carlosonunez/flight-summarizer/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

type MockBrowser struct {
	doc *html.Node
}

func (b *MockBrowser) Init(fNum string) error {
	var err error
	testMap := map[string]string{
		"FAKE1": "live_ontime_departure_early_arrival",
	}
	mock, ok := testMap[fNum]
	if !ok {
		panic(fmt.Sprintf("invalid flightera test case: %s", fNum))
	}
	data, err := testhelpers.LoadFixture(mock)
	if err != nil {
		return err
	}
	b.doc, err = html.Parse(bytes.NewReader(data))
	return err
}

func (b *MockBrowser) Visit(types.BrowserOpts) error {
	return nil
}

func (b *MockBrowser) Document() *html.Node {
	return b.doc
}

func TestOriginAirport(t *testing.T) {
	b := MockBrowser{}
	err := b.Init("FAKE1")
	require.NoError(t, err)
	iata, err := flighteraGetOriginAirport(&b)
	require.NoError(t, err)
	assert.Equal(t, iata, "CLT")
}

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
	want := testhelpers.MustParseTime("2024-12-15T11:32:00-05:00")
	got, err := flighteraGetScheduledDeparture(&b)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}
