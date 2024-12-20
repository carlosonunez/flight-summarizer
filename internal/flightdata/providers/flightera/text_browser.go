package flightera

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/carlosonunez/flight-summarizer/pkg/browser"
	"golang.org/x/net/html"
)

// FlighteraTextBrowser retrieves Flightera data via plain ol' HTTP requests.
type FlighteraTextBrowser struct {
	url *url.URL
	doc *html.Node
}

func (b *FlighteraTextBrowser) Init(fNo string) error {
	u, err := flighteraURL(fNo)
	if err != nil {
		return err
	}
	b.url = u
	return nil
}

func (b *FlighteraTextBrowser) Visit(opts browser.BrowserOpts) error {
	r, err := http.Get(b.url.String())
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("couldn't get Flightera URL '%s'; wanted '%d', but got '%d'", b.url, http.StatusOK, r.StatusCode)
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	b.doc, err = html.Parse(bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func (b *FlighteraTextBrowser) Document() *html.Node {
	return b.doc
}

func NewFlighteraTextBrowser(fNum string) (*FlighteraTextBrowser, error) {
	b := FlighteraTextBrowser{}
	if err := b.Init(fNum); err != nil {
		return nil, err
	}
	return &b, nil
}
