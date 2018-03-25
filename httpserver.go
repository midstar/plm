package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type HttpServer struct {
	measurement *Measurement
	server      *http.Server
}

func CreateHTTPServer(port int, measurement *Measurement) *HttpServer {
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	server := &HttpServer{measurement: measurement, server: srv}

	http.Handle("/", server)
	return server
}

// ServeHTTP handles incoming HTTP requests
func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Start starts the HTTP server. Stop it using the Stop function.
func (s *HttpServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() shutdown reason: %s", err)
		}
	}()
}

// Stop stops the HTTP server.
func (s *HttpServer) Stop() {
	s.server.Shutdown(context.Background())
}
