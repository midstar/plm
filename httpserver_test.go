package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/midstar/proci"
)

func respToString(response io.ReadCloser) string {
	defer response.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(response)
	return buf.String()
}

func TestHttpServer(t *testing.T) {

	// We need to clear default serve mux if http handler is called
	// more than once. We do run it serveral times in the unit tests.
	http.DefaultServeMux = new(http.ServeMux)

	port := 9090
	baseURL := fmt.Sprintf("http://localhost:%d", port)

	// Creata a Measurement object and generate some data
	pMock := proci.GenerateMock(10)
	m := CreateMeasurement(3, 6, 3, 6, pMock)

	m.measureAndLog(false)
	time.Sleep(2 * time.Second) // To make time differ
	timeStamp1 := time.Now()
	m.measureAndLog(true)
	time.Sleep(2 * time.Second) // To make time differ
	timeStamp2 := time.Now()
	time.Sleep(2 * time.Second) // To make time differ
	m.measureAndLog(false)

	// Create and start the HTTP server
	httpServer := CreateHTTPServer("", port, m)
	t.Log("Starting HTTP server")
	httpServer.Start()

	// Tests
	testGetIndex(t, baseURL)
	testGetAllProcesses(t, baseURL)
	testGetProcessesWithUID(t, baseURL)
	testGetProcessesWithMatch(t, baseURL)
	testGetProcessesWithMultipleMatch(t, baseURL)
	testGetAllRAM(t, baseURL)
	testGetPlot(t, baseURL)
	testGetPlotBetween(t, baseURL, timeStamp1, timeStamp2)
	testGetMeasurements(t, baseURL)
	testGetMeasurementsBetween(t, baseURL, timeStamp1, timeStamp2)
	testGetMinMaxMem(t, baseURL)
	testInvalidPath(t, baseURL)

	// Stop HTTP server
	httpServer.Stop()
}

// Called from TestHttpServer
func testGetIndex(t *testing.T, baseURL string) {
	resp, err := http.Get(baseURL)
	if err != nil {
		t.Fatal("Index page not loaded. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	respString := respToString(resp.Body)
	if !strings.Contains(respString, "<title>Process Load Monitor</title>") {
		t.Fatal("Index html title missing")
	}
	if !strings.Contains(respString, "path_5") {
		t.Fatal("Processes are missing")
	}
}

// Called from TestHttpServer
func testGetAllProcesses(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/processes", baseURL))
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var processes map[int]*Process
	err = json.Unmarshal(body, &processes)
	if err != nil {
		t.Fatal("Unable decode get processes. Reason: ", err)
	}
	if len(processes) != 10 {
		t.Fatal("Not all processes received")
	}
}

// Called from TestHttpServer
func testGetProcessesWithUID(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/processes?uids=3,6", baseURL))
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var processes map[int]*Process
	err = json.Unmarshal(body, &processes)
	if err != nil {
		t.Fatal("Unable decode get processes. Reason: ", err)
	}
	if _, hasElement := processes[3]; !hasElement {
		t.Fatal("Process with UID 3 not returned")
	}
	if _, hasElement := processes[6]; !hasElement {
		t.Fatal("Process with UID 6 not returned")
	}
	if len(processes) != 2 {
		t.Fatal("Only two processes expected but got: ", len(processes))
	}

	// Test invalid UID
	resp, err = http.Get(fmt.Sprintf("%s/processes?uids=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetProcessesWithMatch(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/processes?match=path_8", baseURL))
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var processes map[int]*Process
	err = json.Unmarshal(body, &processes)
	if err != nil {
		t.Fatal("Unable decode get processes. Reason: ", err)
	}
	if len(processes) != 1 {
		t.Fatal("Only one process expected but got: ", len(processes))
	}
	for _, process := range processes {
		// Should only be one element
		if process.Pid != 8 {
			t.Fatal("Process with PID 8 expected")
		}
	}
}

// Called from TestHttpServer
func testGetProcessesWithMultipleMatch(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/processes?match=path_8&match=path", baseURL))
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var processes map[int]*Process
	err = json.Unmarshal(body, &processes)
	if err != nil {
		t.Fatal("Unable decode get processes. Reason: ", err)
	}
	if len(processes) != 10 {
		t.Fatal("Not all processes received")
	}
}

// Called from TestHttpServer
func testGetAllRAM(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/ram", baseURL))
	if err != nil {
		t.Fatal("Unable to get ram. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get ram. Reason: ", err)
	}
	var phys PhysicalMemory
	err = json.Unmarshal(body, &phys)
	if err != nil {
		t.Fatal("Unable decode get phys. Reason: ", err)
	}
	if phys.TotalPhys != 4*1024*1024 {
		t.Fatal("Wrong total phys memory")
	}
}

// Called from TestHttpServer
func testInvalidPath(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/invalidpath", baseURL))
	if err != nil {
		t.Fatal("Unable to get invalidpath. Reason: ", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetPlot(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/plot", baseURL))
	if err != nil {
		t.Fatal("Plot page not loaded. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	respString := respToString(resp.Body)
	if !strings.Contains(respString, "<div id=\"plotarea\" style=") {
		t.Fatal("Plot html no plot area")
	}
	if !strings.Contains(respString, "path_5") {
		t.Fatal("Processes are missing")
	}

	// Test invalid UID
	resp, err = http.Get(fmt.Sprintf("%s/plot?uids=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get plot. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetPlotBetween(t *testing.T, baseURL string, from time.Time, to time.Time) {
	queryParams := url.Values{}
	queryParams.Add("from", from.Format(time.RFC3339))
	t.Log("From: ", queryParams.Get("from"))
	queryParams.Add("to", to.Format(time.RFC3339))
	t.Log("To: ", queryParams.Get("to"))
	path := fmt.Sprintf("%s/plot?%s", baseURL, queryParams.Encode())
	resp, err := http.Get(path)
	if err != nil {
		t.Fatal("Plot page not loaded. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	respString := respToString(resp.Body)
	if !strings.Contains(respString, "<div id=\"plotarea\" style=") {
		t.Fatal("Plot html no plot area")
	}
	if !strings.Contains(respString, "path_5") {
		t.Fatal("Processes are missing")
	}

	// Test invalid from
	resp, err = http.Get(fmt.Sprintf("%s/plot?from=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get plot. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetMeasurements(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/measurements", baseURL))
	if err != nil {
		t.Fatal("Unable to get measurements. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var measurements ProcessMeasurements
	err = json.Unmarshal(body, &measurements)
	if err != nil {
		t.Fatal("Unable decode get measurements. Reason: ", err)
	}
	if len(measurements.Memory) != 10 {
		t.Fatal("Expected measurements for all process but got only: ", len(measurements.Memory))
	}
	if len(measurements.Memory[1]) != 3 {
		t.Fatal("Expected 3 measurements but got: ", len(measurements.Memory[1]))
	}
	if len(measurements.Times) != 3 {
		t.Fatal("Expected 3 time stamps but got: ", len(measurements.Times))
	}

	// Test invalid UID
	resp, err = http.Get(fmt.Sprintf("%s/measurements?uids=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get measurements. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetMeasurementsBetween(t *testing.T, baseURL string, from time.Time, to time.Time) {
	queryParams := url.Values{}
	queryParams.Add("from", from.Format(time.RFC3339))
	t.Log("From: ", queryParams.Get("from"))
	queryParams.Add("to", to.Format(time.RFC3339))
	t.Log("To: ", queryParams.Get("to"))
	path := fmt.Sprintf("%s/measurements?%s", baseURL, queryParams.Encode())
	resp, err := http.Get(path)
	if err != nil {
		t.Fatal("Unable to get measurements. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	var measurements ProcessMeasurements
	err = json.Unmarshal(body, &measurements)
	if err != nil {
		t.Fatal("Unable decode get measurements. Reason: ", err)
	}
	for i, timeStamp := range measurements.Times {
		t.Log("Time ", i, ": ", timeStamp)
	}
	if len(measurements.Memory) != 10 {
		t.Fatal("Expected measurements for all process but got only: ", len(measurements.Memory))
	}
	if len(measurements.Memory[1]) != 1 {
		t.Fatal("Expected 1 measurements but got: ", len(measurements.Memory[1]))
	}
	if len(measurements.Times) != 1 {
		t.Fatal("Expected 1 time stamps but got: ", len(measurements.Times))
	}

	// Test invalid to
	resp, err = http.Get(fmt.Sprintf("%s/measurements?to=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get measurements. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}

// Called from TestHttpServer
func testGetMinMaxMem(t *testing.T, baseURL string) {
	resp, err := http.Get(fmt.Sprintf("%s/minmaxmem", baseURL))
	if err != nil {
		t.Fatal("Unable to get minmaxmem. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Unable to get processes. Reason: ", err)
	}
	type ProcessMinMaxMem struct {
		Process
		MaxMemoryInPeriod uint32 // Maximum memory during period (KB)
		MinMemoryInPeriod uint32 // Minimum memory during period(KB)
	}
	var processesMinMaxSlice []ProcessMinMaxMem
	err = json.Unmarshal(body, &processesMinMaxSlice)
	if err != nil {
		t.Fatal("Unable decode get minmaxmem. Reason: ", err)
	}
	if len(processesMinMaxSlice) != 10 {
		t.Fatal("Expected min max for all process but got only: ", len(processesMinMaxSlice))
	}

	// Test invalid UID
	resp, err = http.Get(fmt.Sprintf("%s/minmaxmem?uids=invalid", baseURL))
	if err != nil {
		t.Fatal("Unable to get minmaxmem. Reason: ", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}
}
