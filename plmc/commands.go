package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func processFilters() string {
	queryParams := url.Values{}
	if Matcher != "" {
		queryParams.Add("match", Matcher)
	}
	if UIDs != "" {
		queryParams.Add("uids", UIDs)
	}
	if len(queryParams) == 0 {
		return ""
	}
	return fmt.Sprintf("?%s", queryParams.Encode())
}

// CmdPlot get plot for one or more processes
func CmdPlot(filename string) error {
	resp, err := http.Get(fmt.Sprintf("%s/plot%s", PLMUrl, processFilters()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		return err
	}
	fmt.Println(filename, " written")
	return nil
}

// Process represent one unique process
type Process struct {
	UID           int       // Unique ID
	Pid           uint32    // Process PID
	IsAlive       bool      // Is process alive?
	Path          string    // The process path (and name)
	Name          string    // Name of the process (last part of Path)
	CommandLine   string    // The process command line
	MaxMemoryEver uint32    // Maximum memory ever measured (KB)
	MinMemoryEver uint32    // Minimum memory ever measured (KB)
	LastMemory    uint32    // Last memory measured (KB)
	Created       time.Time // When this process was created (or first seen)
	Died          time.Time // When this process died
}

// CmdInfo list info about for one or more processes
func CmdInfo() error {
	resp, err := http.Get(fmt.Sprintf("%s/processes%s", PLMUrl, processFilters()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var processes map[int]*Process
	err = json.Unmarshal(body, &processes)
	if err != nil {
		return err
	}
	fmt.Println("Number of processes: ", len(processes))
	fmt.Println("")
	for _, process := range processes {
		fmt.Println("-----------------------------------------------")
		fmt.Println("PID:          ", process.Pid)
		fmt.Println("UID:          ", process.UID)
		fmt.Println("Name:         ", process.Name)
		fmt.Println("Path:         ", process.Path)
		fmt.Println("Command line: ", process.CommandLine)
		fmt.Println("Max memory:   ", process.MaxMemoryEver, "KB")
		fmt.Println("Min memory:   ", process.MinMemoryEver, "KB")
		fmt.Println("Last memory:  ", process.LastMemory, "KB")
		fmt.Println("First seen:   ", process.Created)
		fmt.Println("Is alive:     ", process.IsAlive)
		if !process.IsAlive {
			fmt.Println("Died:         ", process.Died)
		}
		fmt.Println("")
	}
	return nil
}
