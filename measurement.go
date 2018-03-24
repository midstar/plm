package main

import (
	"github.com/midstar/proci"
	"sync"
	"time"
)

// Measurement holds all measurements.
type Measurement struct {
	FastLogger *Logger
	SlowLogger *Logger
	PM         *ProcessMap
	Mutex      *sync.Mutex // Only access this struct using this mutex
}

// CreateMeasurement creates a new measurment object
//
// Parameter mutex is used to protect the Measurement struct in different
// threads.
func CreateMeasurement(fastLoggerSize int, slowLoggerSize int,
	mutex *sync.Mutex, pi proci.Interface) *Measurement {
	return &Measurement{
		FastLogger: CreateLogger(fastLoggerSize),
		SlowLogger: CreateLogger(slowLoggerSize),
		PM:         NewProcessMap(pi),
		Mutex:      mutex}
}

// MeasureLoop runs the measurement loop. Supposed to be runned as a goroutine.
//
// Parameter halt is a channel to which you can send a bool value to halt the
// measurement loop (i.e. exit this function).
func (m *Measurement) MeasureLoop(fastLogTimeMs int, slowLogFactor int, halt chan bool) {
	haltMeasurement := false
	addToSlowLog := false
	iter := 1
	for !haltMeasurement {

		if iter%slowLogFactor == 0 {
			addToSlowLog = true
			iter = 0
		} else {
			addToSlowLog = false
		}

		m.measureAndLog(addToSlowLog)

		iter++

		select {
		case <-halt:
			haltMeasurement = true
		case <-time.After(time.Duration(fastLogTimeMs) * time.Millisecond):
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
