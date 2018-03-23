package main

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := CreateLogger(3)
	
	logger.AddRow(&LogRow{
		Time : time.Now(),
		MemUsed : 0,
		LogProcesses : make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 1, logger.NbrRows)
	assertEqualsInt(t, "Next index", 1, logger.Index)
	
	logger.AddRow(&LogRow{
		Time : time.Now(),
		MemUsed : 1,
		LogProcesses : make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 2, logger.NbrRows)
	assertEqualsInt(t, "Next index", 2, logger.Index)
	
	logger.AddRow(&LogRow{
		Time : time.Now(),
		MemUsed : 2,
		LogProcesses : make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 3, logger.NbrRows)
	assertEqualsInt(t, "Next index", 0, logger.Index)
	
	logger.AddRow(&LogRow{
		Time : time.Now(),
		MemUsed : 3,
		LogProcesses : make([]*LogProcess, 0)})
	assertEqualsInt(t, "Number of elements", 3, logger.NbrRows)
	assertEqualsInt(t, "Next index", 1, logger.Index)
	assertEqualsInt(t, "First index MemUsed",  3, int(logger.LogRows[0].MemUsed))
	assertEqualsInt(t, "Second index MemUsed", 1, int(logger.LogRows[1].MemUsed))
	assertEqualsInt(t, "Third index MemUsed",  2, int(logger.LogRows[2].MemUsed))
}