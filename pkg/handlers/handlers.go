package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ParkHandler handles HTTP requests for the main park simulator
type ParkHandler struct{}

// NewParkHandler creates a new park handler
func NewParkHandler() *ParkHandler {
	return &ParkHandler{}
}

// ServeHTTP implements the http.Handler interface
func (h *ParkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/metrics":
		promhttp.Handler().ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}

// AttractionHandler handles HTTP requests for attractions
type AttractionHandler struct{}

// NewAttractionHandler creates a new attraction handler
func NewAttractionHandler() *AttractionHandler {
	return &AttractionHandler{}
}

// ServeHTTP implements the http.Handler interface
func (h *AttractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/metrics":
		promhttp.Handler().ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
