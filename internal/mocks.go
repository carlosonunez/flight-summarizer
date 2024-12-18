package internal

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/carlosonunez/flight-summarizer/types"
	"github.com/gocarina/gocsv"
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

type MockTimeZoneDB struct {
	entries []*types.TimeZoneDBInfo
}

func (d *MockTimeZoneDB) Init() error {
	if err := gocsv.Unmarshal(filepath.Join(testhelpers.FIXTURE_PATH, "timezonedb.csv"), &d.entries); err != nil {
		return err
	}
	return nil
}

// 2024-12-18: TODO: The timezone DB has >124k entries in it; O(n) linear search
// will be slow. At the same time, I don't want to use a separate database (that
// I have to deploy and pay for.
// Find an in-memory database thing that can read CSV or SQL dumps and use gosql
// to query it.
func (d *MockTimeZoneDB) LookupUTCOffsetByID(ID string, start time.Time) (int64, error) {
	return 0, errors.New("not implemented yet")
}
