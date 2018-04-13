package main

import (
	"log"

	"github.com/midstar/proci"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// PLM the PLM context
type PLM struct {
	Config      *Configuration
	httpServer  *HTTPServer
	measurement *Measurement
}

// CreatePLM loads the configuration and creates the HTTP server and
// measurement
func CreatePLM() *PLM {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "plm.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
	log.Print("Startup of PLM")
	configuration := LoadConfiguration(DefaultConfigFile)
	m := CreateMeasurement(configuration.FastLogSize, configuration.SlowLogSize,
		configuration.FastLogTimeMs, configuration.SlowLogSize, proci.Proci{})
	s := CreateHTTPServer(configuration.Port, m)
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

