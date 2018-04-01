package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/midstar/proci"
)

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

// PhysicalMemory represents the physical RAM memory
type PhysicalMemory struct {
	TotalPhys   uint32 // Total memory installed (KB)
	MaxPhysEver uint32 // Maximum used physical memory ever measured (KB)
	MinPhysEver uint32 // Minimum used physical memory ever measured (KB)
	LastPhys    uint32 // Last used physical memory measured (KB)
}

// ProcessMap has two internal maps. Both maps are pointing to the
// same Process objects, but keyed on different identities.
// The reason for this design is that the PID's might be reused by
// the operating system.
type ProcessMap struct {
	nextUniqueID int                 // Increment for each created process
	All          map[int]*Process    // A map with all processes, keyed on UID
	Alive        map[uint32]*Process // A map with the living processes, keyd on PID
	Phys         *PhysicalMemory     // Represents the physical memory
	LastUpdate   time.Time           // Last time this map was updated
	Pi           proci.Interface     // Interface for reading processes
}

// NewProcessMap creates a new process map
func NewProcessMap(pi proci.Interface) *ProcessMap {
	// set only specific field value with field key
	return &ProcessMap{
		nextUniqueID: 0,
		All:          make(map[int]*Process),
		Alive:        make(map[uint32]*Process),
		Phys:         &PhysicalMemory{},
		Pi:           pi}
}

// CreateProcess creates a new process in the ProcessMap. It will assign it
// a unique identity and put it in both the All and Alive maps. If another
// process with the same PID exist in the All map, the old process will be
// set to Alive = false and removed from the All list.
func (processMap *ProcessMap) CreateProcess(pid uint32, fullPath string, commandLine string) *Process {
	processMap.nextUniqueID++
	uid := processMap.nextUniqueID
	_, name := filepath.Split(fullPath)
	process := Process{
		UID:         uid,
		Pid:         pid,
		IsAlive:     true,
		Path:        fullPath,
		Name:        name,
		CommandLine: commandLine,
		Created:     time.Now()}
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
	process.Died = time.Now()
	delete(processMap.Alive, pid)
}

// Update starts with setting living processes to IsAlive = false, then it will
// go through all processes reported by the operating system and update the
// corresponding process in the dictionary. If if a new process is detected a
// new entry in the process map is added.
//
// The Pid, Path, Name and CommandLine fields of the process are only
//updated if the process is new.
//
// If the process is dead it will be removed from the Alive field in ProcessMap.
func (processMap *ProcessMap) Update() {

	// Start with setting all processes to IsAlive = false
	for _, process := range processMap.Alive {
		process.IsAlive = false
	}

	// List and update or create all processes

	pids := processMap.Pi.GetProcessPids()
	for i := 0; i < len(pids); i++ {
		pid := pids[i]
		process, hasPid := processMap.Alive[pid]
		fullPath, patherr := processMap.Pi.GetProcessPath(pid)
		if patherr != nil || fullPath == "" {
			// This is probably a system process that we cannot access.
			// Pointless to track this process
			if hasPid {
				// It was a valid process with pid this before. I must have
				// died
				processMap.ProcessKilled(pid)
			}
			continue
		}
		if hasPid && fullPath != process.Path {
			// The fullPath has changed. It must be a new process that has replaced
			// the old one.
			processMap.ProcessKilled(pid)
			hasPid = false
		}
		if !hasPid {
			// We have a new process
			commandLine, cmderr := processMap.Pi.GetProcessCommandLine(pid)
			if cmderr != nil {
				// Expected for some system processes.
				commandLine = ""
			}
			process = processMap.CreateProcess(pid, fullPath, commandLine)
		}
		process.IsAlive = true

		memoryUsage, memerr := processMap.Pi.GetProcessMemoryUsage(pid)
		if memerr != nil {
			fmt.Println("GetProcessMemoryUsage for PID", pid, "returned error:", memerr)
			process.LastMemory = 0
		} else {
			memoryUsageKB := uint32(memoryUsage / 1024) // Byte to KiloByte
			if process.MinMemoryEver == 0 || memoryUsageKB < process.MinMemoryEver {
				process.MinMemoryEver = memoryUsageKB
			}
			if memoryUsageKB > process.MaxMemoryEver {
				process.MaxMemoryEver = memoryUsageKB
			}
			process.LastMemory = memoryUsageKB
		}
	}

	processMap.LastUpdate = time.Now()

	// Mark all processes not listed as killed
	for pid, process := range processMap.Alive {
		if process.IsAlive == false {
			processMap.ProcessKilled(pid)
		}
	}

	processMap.updatePhysicalMemory()
}

func (processMap *ProcessMap) updatePhysicalMemory() {
	// Update the overall (physical memory)
	memoryStatus, memstaterr := processMap.Pi.GetMemoryStatus()
	if memstaterr != nil {
		fmt.Println("GetMemoryStatus returned error:", memstaterr)
		processMap.Phys.LastPhys = 0
	} else {
		processMap.Phys.TotalPhys = uint32(memoryStatus.TotalPhys / 1024)                          // Byte to KiloByte
		processMap.Phys.LastPhys = processMap.Phys.TotalPhys - uint32(memoryStatus.AvailPhys/1024) // Byte to KiloByte
		if processMap.Phys.LastPhys > processMap.Phys.MaxPhysEver {
			processMap.Phys.MaxPhysEver = processMap.Phys.LastPhys
		}
		if processMap.Phys.MinPhysEver == 0 || processMap.Phys.LastPhys < processMap.Phys.MinPhysEver {
			processMap.Phys.MinPhysEver = processMap.Phys.LastPhys
		}
	}
}
