package main

import (
	"context"
	"github.com/midstar/proci"
	"sync"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	// Creata a Measurement object and collect some data
	mutex := sync.Mutex{}
	m := CreateMeasurement(20, 20, &mutex, proci.Proci{})

	halt := make(chan bool)

	go m.MeasureLoop(200, 2, halt)

	time.Sleep(2 * time.Second)

	// Halt the measurement loop
	halt <- true

	// Start the HTTP server
	t.Log("Starting HTTP server")
	server := StartHTTPServer(9090, m)

	time.Sleep(3 * time.Second)
	server.Shutdown(context.Background())
}
