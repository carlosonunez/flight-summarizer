package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/carlosonunez/flight-summarizer/internal/summarizers"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
)

// SummarizerServerOptions describes options available for creating summarizer
// servers.
type SummarizerServerOptions struct {
	// SummarizerName is the name of the summarizer to use.
	SummarizerName string

	// FlightNumber is the flight number or flight number hint to retrieve a summary for.
	FlightNumber string
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
