package main

import (
	"fmt"
  "net/http"
  "log"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  configuration := loadConfiguration()
  log.Printf("Listening to port: %d", configuration.Port)
  portStr := fmt.Sprintf(":%d", configuration.Port)
  http.HandleFunc("/", handler)
  log.Fatal(http.ListenAndServe(portStr, nil))
}