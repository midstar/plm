package main

import (
	"fmt"
	"github.com/midstar/proci"
)

// Process represent one unique process
type Process struct {
	UID           int    // Unique ID
	Pid           uint32 // Process PID
	IsAlive       bool   // Is process alive?
	Path          string // The process path (and name)
	CommandLine   string // The process command line
	MaxMemoryEver uint32 // Maximum memory ever measured (KB)
	MinMemoryEver uint32 // Minimum memory ever measured (KB)
	CurrMaxMemory uint32 // Current max memory measured (KB) (will be zeroed)
}

// ProcessMap has two internal maps. Both maps are pointing to the
// same Process objects, but keyed on different identities.
// The reason for this design is that the PID's might be reused by
// the operating system.
type ProcessMap struct {
	nextUniqueID int                 // Increment for each created process
	All          map[int]*Process    // A map with all processes, keyed on UID
	Alive        map[uint32]*Process // A map with the living processes, keyd on PID
	Pi           proci.Interface     // Interface for reading processes
}

// NewProcessMap creates a new process map
func NewProcessMap(pi proci.Interface) *ProcessMap {
	// set only specific field value with field key
	return &ProcessMap{
		nextUniqueID: 0,
		All:          make(map[int]*Process),
		Alive:        make(map[uint32]*Process),
		Pi:           pi}
}

// CreateProcess creates a new process in the ProcessMap. It will assign it
// a unique identity and put it in both the All and Alive maps. If another
// process with the same PID exist in the All map, the old process will be
// set to Alive = false and removed from the All list.
func (processMap *ProcessMap) CreateProcess(pid uint32, path string, commandLine string) *Process {
	processMap.nextUniqueID++
	uid := processMap.nextUniqueID
	process := Process{UID: uid, Pid: pid, IsAlive: true, Path: path, CommandLine: commandLine}
	_, hasPid := processMap.Alive[pid]
	if hasPid {
		processMap.ProcessKilled(pid)
	}
	processMap.All[uid] = &process
	processMap.Alive[pid] = &process
	return &process
}

// ProcessKilled removed process from Alive
func (processMap *ProcessMap) ProcessKilled(pid uint32) {
	process := processMap.Alive[pid]
	process.IsAlive = false
	delete(processMap.Alive, pid)
}

// Update starts with setting living processes to IsAlive = false, then it will
// go through all processes reported by the operating system and update the
// corresponding process in the dictionary. If if a new process is detected a
// new entry in the process map is added.
//
// The Pid, Path and CommandLine fields of the process are only updated if the
// process is new.
//
// If the process is dead it will be removed from the Alive field in ProcessMap.
func (processMap *ProcessMap) Update() {

	// Start with setting all processes to IsAlive = false
	for _, process := range processMap.Alive {
		process.IsAlive = false
	}

	// List and update or create all processes

	pids := proci.GetProcessPids()
	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		if pid == 0 {
			// This is the idle process. No operations can be performed
			// on it.
			continue
		}
		process, hasPid := processMap.Alive[pid]
		path, patherr := proci.GetProcessPath(pid)
		if patherr != nil || path == "" {
			// This is probably a system process that we cannot access.
			// Pointless to track this process
			if hasPid {
				// It was a valid process with pid this before. It's must have
				// died
				processMap.ProcessKilled(pid)
			}
			continue
		}
		if !hasPid {
			// We have a new process
			commandLine, cmderr := proci.GetProcessCommandLine(pid)
			if cmderr != nil {
				// Expected for some system processes.
				commandLine = ""
			}
			process = processMap.CreateProcess(pid, path, commandLine)
		}
		process.IsAlive = true

		memoryUsage, memerr := proci.GetProcessMemoryUsage(pid)
		if memerr != nil {
			fmt.Println("GetProcessMemoryUsage for PID", pid, "returned error:", memerr)
		} else {
			memoryUsageKB := uint32(memoryUsage / 1024) // Byte to KiloByte
			if process.MinMemoryEver == 0 {
				process.MinMemoryEver = memoryUsageKB
			}
			if memoryUsageKB > process.MaxMemoryEver {
				process.MaxMemoryEver = memoryUsageKB
			}
			if memoryUsageKB > process.CurrMaxMemory {
				process.CurrMaxMemory = memoryUsageKB
			}
		}
	}

	// Mark all processes not listed as killed
	for pid, process := range processMap.Alive {
		if process.IsAlive == false {
			processMap.ProcessKilled(pid)
		}
	}
}
