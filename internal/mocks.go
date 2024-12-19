package internal

import (
	"bytes"
	"fmt"

	"github.com/carlosonunez/flight-summarizer/testhelpers"
	"github.com/carlosonunez/flight-summarizer/types"
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
