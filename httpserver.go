package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	measurement *Measurement
	server      *http.Server
}

// CreateHTTPServer creates the HTTP server. Start it with Start.
func CreateHTTPServer(port int, measurement *Measurement) *HTTPServer {
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	server := &HTTPServer{measurement: measurement, server: srv}

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
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s!", r.URL.Path)
	}
}

func (s *HTTPServer) serveHTTPIndex(w http.ResponseWriter) {
	funcMap := template.FuncMap{
		// Convert KB to MB only keep one decimal
		"kb_to_mb": func(kb uint32) string {
			return fmt.Sprintf("%.1f", float64(kb)/1024.0)
		},
	}
	t, err := template.New("").Funcs(funcMap).ParseFiles("templates/index.gohtml")
	if err != nil {
		http.Error(w, "Create template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "index.gohtml", s.measurement.PM)
	if err != nil {
		http.Error(w, "Execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) serveHTTPListAllProcesses(w http.ResponseWriter) {
	s.measurement.Mutex.Lock()
	js, err := json.Marshal(s.measurement.PM.All)
	s.measurement.Mutex.Unlock()
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
