package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	measurement *Measurement
	server      *http.Server
	fm          *template.FuncMap
	basePath    string
	tags        map[string]time.Time // Use tagsMutext for read/write
	tagsMutex   sync.Mutex
}

// CreateHTTPServer creates the HTTP server. Start it with Start.
func CreateHTTPServer(basePath string, port int, measurement *Measurement) *HTTPServer {
	funcMap := &template.FuncMap{
		// Convert KB to MB only keep one decimal
		"kb_to_mb": func(kb uint32) string {
			return fmt.Sprintf("%.1f", float64(kb)/1024.0)
		},

		// Convert KB to MB only keep one decimal
		"int_to_str": func(v int) string {
			return fmt.Sprintf("%d", v)
		},

		// Convert an array of kily bytes to megabytes. Keep one decimal.
		"slice_kb_to_mb": func(kb_values []uint32) []float64 {
			mbValues := make([]float64, len(kb_values))
			for i := 0; i < len(kb_values); i++ {
				mbValues[i] = float64(kb_values[i]/512) / 2
			}
			return mbValues
		},

		// Log utulization in %
		"log_utilization": func(log Logger) int {
			return int(float64(log.NbrRows*100) / float64(log.MaxRows))
		}}
	portStr := fmt.Sprintf(":%d", port)
	srv := &http.Server{Addr: portStr}
	server := &HTTPServer{
		basePath:    basePath,
		measurement: measurement,
		server:      srv,
		fm:          funcMap,
		tags:        make(map[string]time.Time),
		tagsMutex:   sync.Mutex{}}

	http.Handle("/", server)
	return server
}

// ServeHTTP handles incoming HTTP requests
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	req := fmt.Sprintf("%s %s", r.Method, segments[1])
	switch req {
	case "GET ":
		s.serveHTTPIndex(w)
	case "GET processes":
		s.serveHTTPListProcesses(w, r.URL.Query())
	case "GET ram":
		s.serveHTTPGetRAM(w)
	case "GET plot":
		s.serveHTTPPlot(w, r.URL.Query())
	case "GET measurements":
		s.serveHTTPMeasurements(w, r.URL.Query())
	case "GET minmaxmem":
		s.serveHTTPGetMinMaxMem(w, r.URL.Query())
	case "POST tag":
		if len(segments) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Now tag name provided")
		} else {
			s.serveHTTPPostTag(w, segments[2])
		}
	case "GET tag":
		if len(segments) < 3 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Now tag name provided")
		} else {
			s.serveHTTPGetTag(w, segments[2])
		}
	case "GET tags":
		s.serveHTTPGetTags(w)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not a valid path: %s or method %s!", r.URL.Path, r.Method)
	}
}

// getFromTo is a function that identifies if the URL includes any of following
// query parameters:
//  - from (restrict result from time in RFC3339 format)
//  - to (restruct result to time in RFC3339 format)
//  - fromTag (as from but use a tag)
//  - toTag (as to but use tag)
//
// If the above is not given the zero (default) time is returned.
//
// Returns:
// (from, to, error)
func (s *HTTPServer) getFromTo(values url.Values) (time.Time, time.Time, error) {
	from := time.Time{} // Zero
	to := time.Time{}   // Zero
	var err error

	fromStr, hasElement := values["from"]
	if hasElement {
		from, err = time.Parse(time.RFC3339, fromStr[0])
		if err != nil {
			return from, to, fmt.Errorf("Invalid parameter from %s. Reason: %s", fromStr[0], err)
		}
	} else {
		var fromTagStr []string
		fromTagStr, hasElement = values["fromTag"]
		if hasElement {
			s.tagsMutex.Lock()
			var hasTag bool
			from, hasTag = s.tags[fromTagStr[0]]
			s.tagsMutex.Unlock()
			if !hasTag {
				return from, to, fmt.Errorf("Invalid tag: %s", fromTagStr[0])
			}
		}
	}

	toStr, hasElement := values["to"]
	if hasElement {
		to, err = time.Parse(time.RFC3339, toStr[0])
		if err != nil {
			return from, to, fmt.Errorf("Invalid parameter to %s. Reason: %s", toStr[0], err)
		}
	} else {
		var toTagStr []string
		toTagStr, hasElement = values["toTag"]
		if hasElement {
			s.tagsMutex.Lock()
			var hasTag bool
			to, hasTag = s.tags[toTagStr[0]]
			s.tagsMutex.Unlock()
			if !hasTag {
				return from, to, fmt.Errorf("Invalid tag: %s", toTagStr[0])
			}
		}
	}

	return from, to, nil
}

// getUIDs is a function that identifies if the URL includes any of follwowing
// query parameters:
//  - uids (list of uids, example uids=12,42,1234)
//  - match (match text, example match=myprocess.exe)
//
// If none of the above query pararameters where listed it is assumed that all
// UIDs shall be used.
//
// It is not possible to combine the above parameters.
//
// Returns a list of uids and error.
func (s *HTTPServer) getUIDs(values url.Values) ([]int, error) {
	// First check query parameter
	uids, err := parseQueryUIDs(values)
	if err != nil {
		return uids, err
	}
	if uids != nil {
		return uids, nil
	}

	// If not given, check the match parameter
	uids = s.parseQueryMatch(values)
	if uids != nil {
		return uids, nil
	}

	// If none given, return all uids
	s.measurement.Mutex.Lock()
	defer s.measurement.Mutex.Unlock()
	uids = make([]int, 0, len(s.measurement.PM.All))
	for uid := range s.measurement.PM.All {
		uids = append(uids, uid)
	}

	return uids, nil
}

// parseQueryUIDs parses the uids query parameter and returns a slice of UIDs.
// If the uids parameter is not provided nil will be returned.
func parseQueryUIDs(values url.Values) ([]int, error) {
	uids, hasElement := values["uids"]
	if !hasElement {
		return nil, nil
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

// parseQueryMatch parses the match query parameter and returns a slice of UIDs.
// More than one match parameters might be added and will be interpreted as OR.
// If the match parameter is not provided nil will be returned.
func (s *HTTPServer) parseQueryMatch(values url.Values) []int {
	match, hasElement := values["match"]
	if !hasElement {
		return nil
	}
	s.measurement.Mutex.Lock()
	defer s.measurement.Mutex.Unlock()

	// Filter out processes that match
	uids := make([]int, 0, len(s.measurement.PM.All))
	for _, m := range match {
		uidsTmp := s.measurement.PM.GetUIDs(m)
		for _, uid := range uidsTmp {
			hasUID := false
			for _, uid2 := range uids {
				if uid == uid2 {
					hasUID = true
					break
				}
			}
			if !hasUID {
				uids = append(uids, uid)
			}
		}
	}
	return uids
}

// serveHTTPGetMinMaxMem returns the highest and lowest memory consumption
// during a specific time
func (s *HTTPServer) serveHTTPGetMinMaxMem(w http.ResponseWriter, values url.Values) {
	type ProcessMinMaxMem struct {
		Process
		MaxMemoryInPeriod uint32 // Maximum memory during period (KB)
		MinMemoryInPeriod uint32 // Minimum memory during period(KB)
	}
	uids, err := s.getUIDs(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	from, to, err := s.getFromTo(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	measurements := s.measurement.GetProcessMeasurementsBetween(uids, from, to) // Thread safe
	s.measurement.Mutex.Lock()
	defer s.measurement.Mutex.Unlock()
	result := make([]ProcessMinMaxMem, 0, len(measurements.Memory))
	for uid, values := range measurements.Memory {
		process, hasElement := s.measurement.PM.All[uid]
		if hasElement {
			p := ProcessMinMaxMem{
				Process:           *process,
				MaxMemoryInPeriod: 0,
				MinMemoryInPeriod: 4294967295} // = 2 ^ 32 - 1
			for _, value := range values {
				if value > p.MaxMemoryInPeriod {
					p.MaxMemoryInPeriod = value
				}
				if value < p.MinMemoryInPeriod {
					p.MinMemoryInPeriod = value
				}
			}
			result = append(result, p)
		}
	}
	js, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *HTTPServer) serveHTTPMeasurements(w http.ResponseWriter, values url.Values) {
	uids, err := s.getUIDs(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	from, to, err := s.getFromTo(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	measurements := s.measurement.GetProcessMeasurementsBetween(uids, from, to) // Thread safe
	js, err := json.Marshal(measurements)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *HTTPServer) serveHTTPPlot(w http.ResponseWriter, values url.Values) {
	uids, err := s.getUIDs(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	from, to, err := s.getFromTo(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	templateFile := filepath.Join(s.basePath, "templates", "plot.gohtml")
	t, err := template.New("").Funcs(*s.fm).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Create template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	type MeasAndProcesses struct {
		Measurements *ProcessMeasurements
		Processes    map[int]*Process
	}
	measAndProcesses := MeasAndProcesses{Processes: make(map[int]*Process)}
	measAndProcesses.Measurements = s.measurement.GetProcessMeasurementsBetween(uids, from, to) // Thread safe
	s.measurement.Mutex.Lock()
	for uid := range measAndProcesses.Measurements.Memory {
		measAndProcesses.Processes[uid] = s.measurement.PM.All[uid]
	}
	err = t.ExecuteTemplate(w, "plot.gohtml", measAndProcesses)
	s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, "Execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *HTTPServer) serveHTTPIndex(w http.ResponseWriter) {
	templateFile := filepath.Join(s.basePath, "templates", "index.gohtml")
	t, err := template.New("").Funcs(*s.fm).ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Create template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.measurement.Mutex.Lock()
	err = t.ExecuteTemplate(w, "index.gohtml", s.measurement)
	s.measurement.Mutex.Unlock()
	if err != nil {
		http.Error(w, "Execute template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HTTPServer) serveHTTPListProcesses(w http.ResponseWriter, values url.Values) {

	uids, err := s.getUIDs(values)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.measurement.Mutex.Lock()
	defer s.measurement.Mutex.Unlock()

	processes := make(map[int]*Process)
	for _, uid := range uids {
		process, hasElement := s.measurement.PM.All[uid]
		if hasElement {
			processes[uid] = process
		}
	}

	js, err := json.Marshal(processes)
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

func (s *HTTPServer) serveHTTPPostTag(w http.ResponseWriter, tagName string) {
	s.tagsMutex.Lock()
	defer s.tagsMutex.Unlock()
	s.tags[tagName] = time.Now()
}

func (s *HTTPServer) serveHTTPGetTag(w http.ResponseWriter, tagName string) {
	s.tagsMutex.Lock()
	defer s.tagsMutex.Unlock()
	t, hasTag := s.tags[tagName]
	if !hasTag {
		errText := fmt.Sprintf("Tag %s not found", tagName)
		http.Error(w, errText, http.StatusNotFound)
		return
	}
	js, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *HTTPServer) serveHTTPGetTags(w http.ResponseWriter) {
	s.tagsMutex.Lock()
	defer s.tagsMutex.Unlock()
	js, err := json.Marshal(s.tags)
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
