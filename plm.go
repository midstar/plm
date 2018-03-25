package main

import (
	"fmt"
	"os"
	"signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	fmt.Println("Please press CTRL+C")
	sig <- c
	fmt.Println("Bye bye")
/*	configuration := LoadConfiguration(DefaultConfigFile)
	m := CreateMeasurement(configuration.FastLogSize, configuration.SlowLogSize, 
	                       sync.Mutex{}, proci.Proci{})
	
	log.Printf("Listening to port: %d", configuration.Port)
	portStr := fmt.Sprintf(":%d", configuration.Port)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(portStr, nil))*/
}
