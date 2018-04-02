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

// GetMemUsed returns memory used for a specific process. If process is not
// listed 0 is returned.
func (lr *LogRow) GetMemUsed(uid int) uint32 {
	for _, logProcess := range lr.LogProcesses {
		if logProcess.UID == uid {
			return logProcess.MemUsed
		}
	}
	return 0
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

// OldestIndexGet the index of the oldest entry. -1 if
// no entries exist.
func (l *Logger) OldestIndex() int {
	if l.NbrRows == 0 {
		return -1
	}
	oldestIndex := l.Index
	if l.NbrRows < l.MaxRows {
		// No wrap yet, oldest is the first element
		oldestIndex = 0
	}
	return oldestIndex
}

// OldestDate returns the date for the oldest log entry. If no entry exist
// the current date is returned.
func (l *Logger) OldestDate() time.Time {
	oldestIndex := l.OldestIndex()
	if oldestIndex == -1 {
		return time.Now()
	}
	return l.LogRows[oldestIndex].Time
}
