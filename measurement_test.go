package main

import (
	"github.com/midstar/proci"
	"sync"
	"testing"
	"time"
)

func TestMeasurement(t *testing.T) {
	m := CreateMeasurement(2, 4, proci.Proci{})

	m.measureAndLog(false)
	assertEqualsInt(t, "Size of FastLogger", 1, m.FastLogger.NbrRows)
	assertEqualsInt(t, "Size of SlowLogger", 0, m.SlowLogger.NbrRows)

	m.measureAndLog(true)
	assertEqualsInt(t, "Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t, "Size of SlowLogger", 1, m.SlowLogger.NbrRows)

	m.measureAndLog(false)
	assertEqualsInt(t, "Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t, "Size of SlowLogger", 1, m.SlowLogger.NbrRows)

	m.measureAndLog(true)
	assertEqualsInt(t, "Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t, "Size of SlowLogger", 2, m.SlowLogger.NbrRows)
}

func TestMeasureLoop(t *testing.T) {
	m := CreateMeasurement(2, 4, proci.Proci{})
	mutex := sync.Mutex{}
	halt := make(chan bool)
	
	go m.MeasureLoop(100, 2, &mutex, halt)
	
	time.Sleep(1 * time.Second)
	
	// Halt the measurement loop
	halt <- true
	
}
