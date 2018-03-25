package main

import (
	"github.com/midstar/proci"
	"sync"
	"time"
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
		MemUsed:      m.PM.LastPhys,
		LogProcesses: logProcesses}

	m.FastLogger.AddRow(&row)
	if addToSlowLogger {
		m.SlowLogger.AddRow(&row)
	}

	m.Mutex.Unlock()
}
