package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/carlosonunez/flight-summarizer/internal/summarizers"
	"github.com/carlosonunez/flight-summarizer/pkg/summarizer"
)

// ServerOptions are options to provide to the server.
type ServerOptions struct {
	// Port is the port to bind the summarizer server to (default: 8080)
	Port int

	// ListenAddress specifies the IP address to bind the summarizer to (default:
	// 127.0.0.1)
	ListenAddress string
}

type response struct {
	Status  string                    `json:"status"`
	Error   string                    `json:"error,omitempty"`
	Summary *summarizer.FlightSummary `json:"summary,omitempty"`
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
		return errors.New("server port missing")
	}
	if opts.ListenAddress == "" {
		return errors.New("server address missing")
	}
	mux := http.NewServeMux()
	initRoutes(mux)
	log.Printf("server started at %s:%d; press CTRL-C to stop\n", opts.ListenAddress, opts.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", opts.ListenAddress, opts.Port), mux)
	return nil
}

func initRoutes(mux *http.ServeMux) {
	mux.Handle("/summarize", &handler{fn: summarizeHandler})
	mux.Handle("/ping", &handler{fn: pingHandler})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(&response{Status: "ok"})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ping check failed: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", body)
}

func summarizeHandler(w http.ResponseWriter, r *http.Request) {
	flightNumber, err := lookupFlightNumber(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed while looking up flight: %s", err)
		return
	}
	opts := summarizer.FlightSummarizerOptions{
		FlightNumber: flightNumber,
	}
	s, err := lookupSummarizer(r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "couldn't find summarizer: %s", err)
		return
	}
	if err := initSummarizer(s, &opts); err != nil {
		writeError(w, http.StatusInternalServerError, "failed while initializing this summarizer: %s", err)
		return
	}
	resp, err := summarize(s)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed while summarizing: %s", err)
		return
	}
	writeOK(w, resp)
}

func summarize(s summarizer.FlightSummarizer) (*summarizer.FlightSummary, error) {
	summary, err := s.Summarize()
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func initSummarizer(s summarizer.FlightSummarizer, opts *summarizer.FlightSummarizerOptions) error {
	if err := s.Init(opts); err != nil {
		return err
	}
	return nil
}

func lookupSummarizer(r *http.Request) (summarizer.FlightSummarizer, error) {
	var summarizer summarizer.FlightSummarizer
	var err error
	summarizerStr := r.URL.Query().Get("summarizer")
	if len(summarizerStr) == 0 {
		if summarizer, err = summarizers.LookupDefault(); err != nil {
			return nil, err
		}
	} else {
		if summarizer, err = summarizers.Lookup(summarizerStr); err != nil {
			return nil, err
		}
	}
	return summarizer, nil
}

func lookupFlightNumber(r *http.Request) (string, error) {
	f := r.URL.Query().Get("flightNumber")
	if len(f) == 0 {
		return "", errors.New("flight number missing")
	}
	return f, nil
}

func writeError(w http.ResponseWriter, code int, format string, parts ...any) {
	r := response{
		Status: "error",
		Error:  fmt.Sprintf(format, parts...),
	}
	body, _ := json.Marshal(&r)
	http.Error(w, string(body), code)
}

func writeOK(w http.ResponseWriter, summary *summarizer.FlightSummary) {
	r := response{
		Status:  "ok",
		Summary: summary,
	}
	body, err := json.Marshal(&r)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "couldn't parse summary: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", body)
}
