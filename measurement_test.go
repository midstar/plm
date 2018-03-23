package main

import (
	"testing"
	"github.com/midstar/proci"
)

func TestMeasurement(t *testing.T) {
	m := CreateMeasurement(2, 4, proci.Proci{})
	
	m.MeasureAndLog(false)
	assertEqualsInt(t,"Size of FastLogger", 1, m.FastLogger.NbrRows)
	assertEqualsInt(t,"Size of SlowLogger", 0, m.SlowLogger.NbrRows)
	
	m.MeasureAndLog(true)
	assertEqualsInt(t,"Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t,"Size of SlowLogger", 1, m.SlowLogger.NbrRows)
	
	m.MeasureAndLog(false)
	assertEqualsInt(t,"Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t,"Size of SlowLogger", 1, m.SlowLogger.NbrRows)
	
	m.MeasureAndLog(true)
	assertEqualsInt(t,"Size of FastLogger", 2, m.FastLogger.NbrRows)
	assertEqualsInt(t,"Size of SlowLogger", 2, m.SlowLogger.NbrRows)
}