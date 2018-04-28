package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func plmPath(t *testing.T) string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		t.Fatal("Environment GOPATH needs to be set for this test")
	}
	return filepath.Join(gopath, "src", "github.com", "midstar", "plm")
}

func TestPLM(t *testing.T) {
	// We need to clear default serve mux if http handler is called
	// more than once. We do run it serveral times in the unit tests.
	http.DefaultServeMux = new(http.ServeMux)

	plm := CreatePLM(plmPath(t))
	plm.Start()
	time.Sleep(3 * time.Second) // Allow some measurements to be done

	resp, err := http.Get("http://localhost:12124")
	if err != nil {
		t.Fatal("Unable to connect to PLM. Reason: ", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("Unexpected status code: ", resp.StatusCode)
	}

	plm.Stop()
}
