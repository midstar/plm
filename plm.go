package main

import (
	"fmt"
	"github.com/midstar/proci"
)

// PLM the PLM context
type PLM struct {
	httpServer  *HTTPServer
	measurement *Measurement
}

// CreatePLM loads the configuration and creates the HTTP server and
// measurement
func CreatePLM() *PLM {
	configuration := LoadConfiguration(DefaultConfigFile)
	m := CreateMeasurement(configuration.FastLogSize, configuration.SlowLogSize,
		configuration.FastLogTimeMs, configuration.SlowLogSize, proci.Proci{})
	s := CreateHTTPServer(configuration.Port, m)
	return &PLM{
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
	fmt.Println("PLM is running. Enter 'exit' to shutdown and exit.")
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
