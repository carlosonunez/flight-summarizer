package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/carlosonunez/flight-summarizer/internal/summarizers"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPort8080               int    = 8080
	defaultListenAddressLocalhost string = "127.0.0.1"
)

// ServerOptions are options to provide to the server.
type ServerOptions struct {
	// Port is the port to bind the summarizer server to (default: 8080)
	Port int

	// ListenAddress specifies the IP address to bind the summarizer to (default:
	// 127.0.0.1)
	ListenAddress string
}

type muxHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	fn func(http.ResponseWriter, *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

// SummarizerServerOptions describes options available for creating summarizer
// servers.
type SummarizerServerOptions struct {
	// SummarizerName is the name of the summarizer to use.
	SummarizerName string

	// FlightNumber is the flight number or flight number hint to retrieve a summary for.
	FlightNumber string
}

// Start starts a new blocking summarizer HTTP server
func Start(opts *ServerOptions) error {
	if opts.Port == 0 {
		log.Warningf("No port provided to the summarizer server; using default port %d", defaultPort8080)
		opts.Port = defaultPort8080
	}
	if opts.ListenAddress == "" {
		log.Warning("No listen address provided; listening on localhost only (external connections might not work)")
		opts.ListenAddress = defaultListenAddressLocalhost
	}
	mux := http.NewServeMux()
	mux.Handle("/summarize", &handler{fn: summarizeHandler})
	http.ListenAndServe(fmt.Sprintf("%s:%d", opts.ListenAddress, opts.Port), mux)
	return nil
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	s, err := lookupSummarizer(w, r)
	if err != nil {
		return
	}
	flightNumber, err := lookupFlightNumber(w, r)
	if err != nil {
		return
	}
	opts := summarizer.FlightSummarizerOptions{
		FlightNumber: flightNumber,
	}
	if err := initSummarizer(s, &opts, w); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	resp, err := summarize(s, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(resp))
}

func summarize(s summarizer.FlightSummarizer, w http.ResponseWriter) ([]byte, error) {
	summary, err := s.Summarize()
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return []byte{}, err
	}
	resp, err := json.Marshal(summary)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return []byte{}, err
	}
	return resp, nil
}

func initSummarizer(s summarizer.FlightSummarizer, opts *summarizer.FlightSummarizerOptions, w http.ResponseWriter) error {
	if err := s.Init(opts); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return err
	}
	return nil
}

func lookupSummarizer(w http.ResponseWriter, r *http.Request) (summarizer.FlightSummarizer, error) {
	var summarizer summarizer.FlightSummarizer
	var err error
	summarizerStr := r.URL.Query().Get("summarizer")
	if len(summarizerStr) == 0 {
		if summarizer, err = summarizers.LookupDefault(); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
			return nil, err
		}
	} else {
		if summarizer, err = summarizers.Lookup(summarizerStr); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return nil, err
		}
	}
	return summarizer, nil
}

func lookupFlightNumber(w http.ResponseWriter, r *http.Request) (string, error) {
	f := r.URL.Query().Get("flightNumber")
	if len(f) == 0 {
		http.Error(w, "flight number missing", http.StatusBadRequest)
		return "", errors.New("flight number missing")
	}
	return f, nil
}
