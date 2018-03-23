package main

import (
	"github.com/midstar/proci"
)

// Measurement holds all measurements.
type Measurement struct {
	FastLogger *Logger
	SlowLogger *Logger
	PM         *ProcessMap
}

// CreateMeasurement creates a new measurment object
func CreateMeasurement(fastLoggerSize int, slowLoggerSize int, pi proci.Interface) *Measurement {
	return &Measurement{
		FastLogger: CreateLogger(fastLoggerSize),
		SlowLogger: CreateLogger(slowLoggerSize),
		PM:         NewProcessMap(pi)}
}

// MeasureAndLog performs measurement and add to FastLogger. Optionally also log to SlowLogger.
func (m *Measurement) MeasureAndLog(addToSlowLogger bool) {

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

}
