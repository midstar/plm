package main

import (
	"github.com/midstar/proci"
	"sync"
	"testing"
	"time"
)

func TestMeasurement(t *testing.T) {
	m := CreateMeasurement(2, 4, &sync.Mutex{}, proci.Proci{})

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
	mutex := sync.Mutex{}
	m := CreateMeasurement(20, 20, &mutex, proci.Proci{})

	halt := make(chan bool)

	go m.MeasureLoop(500, 2, halt)

	time.Sleep(3 * time.Second)

	// Halt the measurement loop
	halt <- true

	t.Log("Size of Fastlogger:", m.FastLogger.NbrRows)
	t.Log("Size of SlowLogger:", m.SlowLogger.NbrRows)
	assertTrue(t, "Size of FastLogger", m.FastLogger.NbrRows > 4 && m.FastLogger.NbrRows < 8)
	assertEqualsInt(t, "Size of SlowLogger", int(m.FastLogger.NbrRows/2), m.SlowLogger.NbrRows)

}

func TestMeasureLoop2(t *testing.T) {
	mutex := sync.Mutex{}
	m := CreateMeasurement(20, 20, &mutex, proci.Proci{})

	halt := make(chan bool)

	go m.MeasureLoop(200, 3, halt)

	time.Sleep(3 * time.Second)

	// Halt the measurement loop
	halt <- true

	t.Log("Size of Fastlogger:", m.FastLogger.NbrRows)
	t.Log("Size of SlowLogger:", m.SlowLogger.NbrRows)
	assertTrue(t, "Size of FastLogger", m.FastLogger.NbrRows > 12 && m.FastLogger.NbrRows < 18)
	assertEqualsInt(t, "Size of SlowLogger", int(m.FastLogger.NbrRows/3), m.SlowLogger.NbrRows)

}
