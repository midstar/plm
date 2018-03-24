package main

import (
	"time"
)

// LogProcess represents one memory measurement for one process
type LogProcess struct {
	UID     int    // Process unique ID (not same as PID, which is not unique)
	MemUsed uint32 // Measured memory used by the process
}

// LogRow represents measurements from all living processes at
// a certain time
type LogRow struct {
	Time         time.Time     // Time when data was measured
	MemUsed      uint32        // Measured total memory used (by all processes)
	LogProcesses []*LogProcess // All process entries
}

// Logger is a collection of LogRows. It is a circular buffer.
type Logger struct {
	LogRows []*LogRow // All rows (circular buffer)
	MaxRows int       // Maximum rows in LogRows
	NbrRows int       // Number of rows written (saturates at MaxRows)
	Index   int       // Next index to write in LogRows
}

// CreateLogger creates a logger with a certain size
func CreateLogger(size int) *Logger {
	return &Logger{
		LogRows: make([]*LogRow, size, size),
		MaxRows: size,
		NbrRows: 0,
		Index:   0}
}

// AddRow adds a new row to the logger
func (l *Logger) AddRow(row *LogRow) {
	l.LogRows[l.Index] = row
	l.Index++
	if l.Index >= l.MaxRows {
		l.Index = 0 // Wrap of log
	}
	if l.NbrRows < l.MaxRows {
		l.NbrRows++
	}
}