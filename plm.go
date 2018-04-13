package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/midstar/proci"
)

// PLM the PLM context
type PLM struct {
	Config      *Configuration
	httpServer  *HTTPServer
	measurement *Measurement
}

// CreatePLM loads the configuration and creates the HTTP server and
// measurement.
//
// basePath is the location where to store log files, read config files and
// templates
func CreatePLM(basePath string) *PLM {

	logFile, err := os.OpenFile(filepath.Join(basePath, "plm.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(mw)

	// Rest of configuration
	log.Print("Startup of PLM")
	configuration := LoadConfiguration(filepath.Join(basePath, DefaultConfigFile))
	log.Print("Listening to port: ", configuration.Port)
	m := CreateMeasurement(configuration.FastLogSize, configuration.SlowLogSize,
		configuration.FastLogTimeMs, configuration.SlowLogSize, proci.Proci{})
	s := CreateHTTPServer(basePath, configuration.Port, m)
	return &PLM{
		Config:      configuration,
		httpServer:  s,
		measurement: m}
}

// Start starts the measurements and HTTP server.
func (plm *PLM) Start() {
	plm.measurement.Start()
	plm.httpServer.Start()
}

// Stop stops the HTTP server and measurement.
func (plm *PLM) Stop() {
	plm.httpServer.Stop()
	plm.measurement.Stop()
}
