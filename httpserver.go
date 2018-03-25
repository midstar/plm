package main

import (
	"fmt"
	"log"
	"net/http"
)

type serverContext struct {
	measurement *Measurement
	server      *http.Server
}

// ServeHTTP handles incoming HTTP requests
func (context *serverContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// StartHTTPServer starts the HTTP server and returns the server instance.
// The server is stopped by calling Shutdown on the server instance.
func StartHTTPServer(port int, measurement *Measurement) *http.Server {
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}

	http.Handle("/", &serverContext{measurement: measurement, server: srv})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() shutdown reason: %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}
