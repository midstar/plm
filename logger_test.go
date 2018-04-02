package main

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := CreateLogger(3)

	currentTime := time.Now()
	time.Sleep(50 * time.Millisecond) // To make time differ
	assertTrue(t, "Oldest date on empty log", currentTime.Before(logger.OldestDate()))
	assertEqualsInt(t, "Oldest index on empty list", -1, logger.OldestIndex())

	logger.AddRow(&LogRow{
		Time:         time.Now(),
		MemUsed:      0,
		LogProcesses: make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 1, logger.NbrRows)
	assertEqualsInt(t, "Next index", 1, logger.Index)
	t.Log("Row 1", logger.LogRows[0].Time)
	assertTrue(t, "Oldest date", logger.LogRows[0].Time == logger.OldestDate())
	assertEqualsInt(t, "Oldest index", 0, logger.OldestIndex())

	time.Sleep(50 * time.Millisecond) // To make time differ
	logger.AddRow(&LogRow{
		Time:         time.Now(),
		MemUsed:      1,
		LogProcesses: make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 2, logger.NbrRows)
	assertEqualsInt(t, "Next index", 2, logger.Index)
	assertTrue(t, "Oldest date", logger.LogRows[0].Time == logger.OldestDate())
	assertEqualsInt(t, "Oldest index", 0, logger.OldestIndex())

	time.Sleep(50 * time.Millisecond) // To make time differ
	logger.AddRow(&LogRow{
		Time:         time.Now(),
		MemUsed:      2,
		LogProcesses: make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 3, logger.NbrRows)
	assertEqualsInt(t, "Next index", 0, logger.Index)
	assertTrue(t, "Oldest date", logger.LogRows[0].Time == logger.OldestDate())
	assertEqualsInt(t, "Oldest index", 0, logger.OldestIndex())

	time.Sleep(50 * time.Millisecond) // To make time differ
	logger.AddRow(&LogRow{
		Time:         time.Now(),
		MemUsed:      3,
		LogProcesses: make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 3, logger.NbrRows)
	assertEqualsInt(t, "Next index", 1, logger.Index)
	assertEqualsInt(t, "First index MemUsed", 3, int(logger.LogRows[0].MemUsed))
	assertEqualsInt(t, "Second index MemUsed", 1, int(logger.LogRows[1].MemUsed))
	assertEqualsInt(t, "Third index MemUsed", 2, int(logger.LogRows[2].MemUsed))
	assertTrue(t, "Oldest date", logger.LogRows[1].Time == logger.OldestDate())
	assertEqualsInt(t, "Oldest index", 1, logger.OldestIndex())

	// Get mem used on non existing process
	assertEqualsInt(t, "Get mem used for non existing process", 0, int(logger.LogRows[0].GetMemUsed(234432)))
}
