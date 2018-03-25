package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	measurement *Measurement
	server      *http.Server
}

// CreateHTTPServer creates the HTTP server. Start it with Start.
func CreateHTTPServer(port int, measurement *Measurement) *HTTPServer {
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	server := &HTTPServer{measurement: measurement, server: srv}

	http.Handle("/", server)
	return server
}

// ServeHTTP handles incoming HTTP requests
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Start starts the HTTP server. Stop it using the Stop function.
func (s *HTTPServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() shutdown reason: %s", err)
		}
	}()
}

// Stop stops the HTTP server.
func (s *HTTPServer) Stop() {
	s.server.Shutdown(context.Background())
}
