package main

import (
	"log"
	"sync"
	"time"

	"github.com/midstar/proci"
)

// Measurement holds all measurements.
type Measurement struct {
	FastLogger    *Logger
	SlowLogger    *Logger
	PM            *ProcessMap
	FastLogTimeMs int
	SlowLogFactor int
	Mutex         *sync.Mutex // Only access this struct using this mutex
	halt          chan bool   // Send to halt measurement
}

// ProcessMeasurements are measuremens from an individual process extracted
// from the Measurement struct. Lengths of all arrays are the same, including
// time. If no measurement was found for a certain time, the measured value
// is set to 0.
type ProcessMeasurements struct {
	Memory map[int][]uint32 // Keyed on UID, values are all measured memory
	Times  []time.Time      // Time values
}

// CreateMeasurement creates a new measurment object
func CreateMeasurement(fastLoggerSize int, slowLoggerSize int,
	fastLogTimeMs int, slowLogFactor int,
	pi proci.Interface) *Measurement {
	return &Measurement{
		FastLogger:    CreateLogger(fastLoggerSize),
		SlowLogger:    CreateLogger(slowLoggerSize),
		PM:            NewProcessMap(pi),
		FastLogTimeMs: fastLogTimeMs,
		SlowLogFactor: slowLogFactor,
		Mutex:         &sync.Mutex{},
		halt:          make(chan bool)}
}

// Start starts the measurement as a separate goroutine.
//
// Stop the measurement with Stop.
func (m *Measurement) Start() {
	go m.measureLoop()
}

// Stop stops the measurement.
func (m *Measurement) Stop() {
	m.halt <- true
}

// GetProcessMeasurementsBetween same as GetProcessMeasurements but only extracts
// measuared values between from and to.
// If from and/or to are set to zero values (default) no restriction is set.
func (m *Measurement) GetProcessMeasurementsBetween(uids []int, from time.Time, to time.Time) *ProcessMeasurements {
	m.Mutex.Lock()
	fastLogOldestTime := m.FastLogger.OldestDate()
	maxSize := m.SlowLogger.NbrRows + m.FastLogger.NbrRows
	pm := &ProcessMeasurements{
		Memory: make(map[int][]uint32),
		Times:  make([]time.Time, 0, maxSize)}
	for _, uid := range uids {
		_, hasElement := m.PM.All[uid]
		if hasElement {
			pm.Memory[uid] = make([]uint32, 0, maxSize)
		} else {
			log.Printf("Trying to get measurement for process with UID %d which don't exist", uid)
		}
	}

	// Start with extracting values from the Slow Log
	slowIndex := m.SlowLogger.OldestIndex()
	handledRows := 0
	for handledRows < m.SlowLogger.NbrRows {
		row := m.SlowLogger.LogRows[slowIndex]
		if row.Time == fastLogOldestTime || row.Time.After(fastLogOldestTime) {
			// Continue with the fast log
			break
		}
		// Only add time if to / from restrictions are fullfilled
		if (from.IsZero() || !row.Time.Before(from)) && (to.IsZero() || !row.Time.After(to)) {
			pm.Times = append(pm.Times, row.Time)
			for uid := range pm.Memory {
				pm.Memory[uid] = append(pm.Memory[uid], row.GetMemUsed(uid))
			}
		}
		handledRows++
		slowIndex++
		if slowIndex == m.SlowLogger.MaxRows {
			// Wrap of log
			slowIndex = 0
		}
	}

	// Continue extract from the fast log
	fastIndex := m.FastLogger.OldestIndex()
	handledRows = 0
	for handledRows < m.FastLogger.NbrRows {
		row := m.FastLogger.LogRows[fastIndex]
		// Only add time if to / from restrictions are fullfilled
		if (from.IsZero() || !row.Time.Before(from)) && (to.IsZero() || !row.Time.After(to)) {
			pm.Times = append(pm.Times, row.Time)
			for uid := range pm.Memory {
				pm.Memory[uid] = append(pm.Memory[uid], row.GetMemUsed(uid))
			}
		}
		handledRows++
		fastIndex++
		if fastIndex == m.FastLogger.MaxRows {
			// Wrap of log
			fastIndex = 0
		}
	}
	m.Mutex.Unlock()
	return pm
}

// GetProcessMeasurements "extracts" the measured values for the provided list
// of processes (using UID as selector) .
func (m *Measurement) GetProcessMeasurements(uids []int) *ProcessMeasurements {
	return m.GetProcessMeasurementsBetween(uids, time.Time{}, time.Time{})
}

// measureLoop runs the measurement loop. Supposed to be runned as a goroutine.
func (m *Measurement) measureLoop() {
	haltMeasurement := false
	addToSlowLog := false
	iter := 1
	for !haltMeasurement {

		if iter%m.SlowLogFactor == 0 {
			addToSlowLog = true
			iter = 0
		} else {
			addToSlowLog = false
		}

		m.measureAndLog(addToSlowLog)
		m.removeOldProcesses()

		iter++

		select {
		case <-m.halt:
			haltMeasurement = true
		case <-time.After(time.Duration(m.FastLogTimeMs) * time.Millisecond):
		}
	}
}

// measureAndLog performs measurement and add to FastLogger. Optionally also log to SlowLogger.
func (m *Measurement) measureAndLog(addToSlowLogger bool) {
	m.Mutex.Lock()

	m.PM.Update()

	logProcesses := make([]*LogProcess, len(m.PM.Alive), len(m.PM.Alive))
	i := 0
	for _, process := range m.PM.Alive {
		logProcesses[i] = &LogProcess{
			UID:     process.UID,
			MemUsed: process.LastMemory}
		i++
	}

	row := LogRow{
		Time:         m.PM.LastUpdate,
		MemUsed:      m.PM.Phys.LastPhys,
		LogProcesses: logProcesses}

	m.FastLogger.AddRow(&row)
	if addToSlowLogger {
		m.SlowLogger.AddRow(&row)
	}

	m.Mutex.Unlock()
}

// removeOldProcesses removes all dead processes where no log entries exists.
func (m *Measurement) removeOldProcesses() {
	m.Mutex.Lock()
	// Only remove entries if the log is full
	if m.SlowLogger.NbrRows == m.SlowLogger.MaxRows {
		oldestTime := m.SlowLogger.OldestDate()
		for uid, process := range m.PM.All {
			if !process.IsAlive && process.Died.Before(oldestTime) {
				log.Printf("Removing process %d. Died: %s, Last entry: %s", uid, process.Died.Format(time.RFC3339), oldestTime.Format(time.RFC3339))
				delete(m.PM.All, uid)
			}
		}
	}
	m.Mutex.Unlock()
}
