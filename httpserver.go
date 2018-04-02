package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	measurement *Measurement
	server      *http.Server
	fm          *template.FuncMap
}

// CreateHTTPServer creates the HTTP server. Start it with Start.
func CreateHTTPServer(port int, measurement *Measurement) *HTTPServer {
	funcMap := &template.FuncMap{
		// Convert KB to MB only keep one decimal
		"kb_to_mb": func(kb uint32) string {
			return fmt.Sprintf("%.1f", float64(kb)/1024.0)
		},
	}
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	server := &HTTPServer{measurement: measurement, server: srv, fm: funcMap}

	http.Handle("/", server)
	return server
}

// ServeHTTP handles incoming HTTP requests
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		//http.ServeFile(w, r, "templates/index.gohtml")
		s.serveHTTPIndex(w)
	case "/processes":
		s.serveHTTPListAllProcesses(w)
	case "/ram":
		s.serveHTTPGetRAM(w)
	case "/plot":
		s.serveHTTPPlot(w, r.URL.Query())
	case "/measurements":
		s.serveHTTPMeasurements(w, r.URL.Query())
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s!", r.URL.Path)
	}
}

// getUIDs parses the uid query parameter and returns a slice of UIDs.
func getUIDs(values url.Values) ([]int, error) {
	uids, hasElement := values["uids"]
	if !hasElement {
		return make([]int, 0), fmt.Errorf("Parameter uids was not provided")
	}
	uidSlice := strings.Split(uids[0], ",")
	uidsInt := make([]int, 0, len(uidSlice))
	for _, uidStr := range uidSlice {
		intValue, valueerr := strconv.Atoi(uidStr)
		if valueerr != nil {
			return make([]int, 0), fmt.Errorf("Invalid parameter uids. UID %s is not a valid integer", uidStr)
		}
		uidsInt = append(uidsInt, intValue)
	}
	return uidsInt, nil
}

func (s *HTTPServer) serveHTTPMeasurements(w http.ResponseWriter, values url.Values) {
	uids, uidsError := getUIDs(values)
	if uidsError != nil {
		http.Error(w, uidsError.Error(), http.StatusBadRequest)
		return
	}
	measurements := s.measurement.GetProcessMeasurements(uids) // Thread safe
	js, err := json.Marshal(measurements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *HTTPServer) serveHTTPPlot(w http.ResponseWriter, values url.Values) {
	uids, uidsError := getUIDs(values)
	if uidsError != nil {
		http.Error(w, uidsError.Error(), http.StatusBadRequest)
		return
	}
	t, err := template.New("").Funcs(*s.fm).ParseFiles("templates/plot.gohtml")
	if err != nil {
		http.Error(w, "Create template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.measurement.Mutex.Lock()
	err = t.ExecuteTemplate(w, "plot.gohtml", uids)
	s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, "Execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *HTTPServer) serveHTTPIndex(w http.ResponseWriter) {
	t, err := template.New("").Funcs(*s.fm).ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, "Create template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.measurement.Mutex.Lock()
	err = t.ExecuteTemplate(w, "index.gohtml", s.measurement.PM)
	s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, "Execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) serveHTTPListAllProcesses(w http.ResponseWriter) {
	s.measurement.Mutex.Lock()
	js, err := json.Marshal(s.measurement.PM.All)
	defer s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *HTTPServer) serveHTTPGetRAM(w http.ResponseWriter) {
	s.measurement.Mutex.Lock()
	js, err := json.Marshal(s.measurement.PM.Phys)
	s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Start starts the HTTP server. Stop it using the Stop function.
func (s *HTTPServer) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() shutdown reason: %s", err)
		}
	}()
}

// Stop stops the HTTP server.
func (s *HTTPServer) Stop() {
	s.server.Shutdown(context.Background())
}
