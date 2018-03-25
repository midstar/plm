package main

import (
	"github.com/midstar/proci"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	// Creata a Measurement object and collect some data
	m := CreateMeasurement(20, 20, 200, 2, proci.Proci{})
	m.StartMeasurement()
	time.Sleep(2 * time.Second)
	m.StopMeasurement()

	// Create and start the HTTP server
	httpServer := CreateHTTPServer(9090, m)
	t.Log("Starting HTTP server")
	httpServer.Start()
	time.Sleep(3 * time.Second)
	httpServer.Stop()
}
