package main

import (
	"fmt"
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

func main() {
	plm := CreatePLM()
	plm.Start()
	fmt.Printf("PLM is running on port %d. Enter 'exit' to shutdown and exit.\n", plm.Config.Port)
	fmt.Print(": ")
	var input string
	for true {
		fmt.Scanln(&input)
		if input == "exit" {
			break
		} else if input == "help" {
			fmt.Println("Supported commands:")
			fmt.Println("  exit : shutdown server and exit")
		} else if input == "" {

		} else {
			fmt.Println("Invalid command. Type 'help' for available commands")
		}
		fmt.Print(": ")
	}
	fmt.Println("Shutting down")
	plm.Stop()
	fmt.Println("Bye bye")
	/*	configuration := LoadConfiguration(DefaultConfigFile)
		m := CreateMeasurement(configuration.FastLogSize, configuration.SlowLogSize,
		                       sync.Mutex{}, proci.Proci{})

		log.Printf("Listening to port: %d", configuration.Port)
		portStr := fmt.Sprintf(":%d", configuration.Port)
		http.HandleFunc("/", handler)
		log.Fatal(http.ListenAndServe(portStr, nil))*/
}
