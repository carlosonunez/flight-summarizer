package flightera

import (
	"bytes"
	"fmt"
	"os"

	"github.com/carlosonunez/flight-summarizer/pkg/browser"
	"golang.org/x/net/html"
)

const (
	LiveOnTimeDepartureEarlyArrival testScenario = iota
	LiveLateDepartureEarlyArrival
)

var testScenarioMap map[testScenario]string = map[testScenario]string{
	LiveOnTimeDepartureEarlyArrival: "fixtures/nojs/live_ontime_departure_early_arrival",
	LiveLateDepartureEarlyArrival:   "fixtures/nojs/live_late_departure_early_arrival",
}

type testScenario int64

type MockBrowser struct {
	doc *html.Node
}

func (b *MockBrowser) Init(mock string) error {
	data, err := os.ReadFile(mock + ".html")
	if err != nil {
		return err
	}
	b.doc, err = html.Parse(bytes.NewReader(data))
	return err
}

func (b *MockBrowser) Visit(browser.BrowserOpts) error {
	return nil
}

func (b *MockBrowser) Document() *html.Node {
	return b.doc
}

func NewMockBrowser(scenario testScenario) (*MockBrowser, error) {
	b := MockBrowser{}
	mock, ok := testScenarioMap[scenario]
	if !ok {
		return &b, fmt.Errorf("no test scenario found that matches ID: %d", scenario)
	}
	if err := b.Init(mock); err != nil {
		return &b, err
	}
	return &b, nil
}
