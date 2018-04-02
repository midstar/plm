package main

import (
	"testing"
	"time"

	"github.com/midstar/proci"
)

func TestGetProcessMeasurements(t *testing.T) {
	pMock := proci.GenerateMock(10)
	m := CreateMeasurement(2, 4, 2, 4, pMock)

	m.measureAndLog(false)
	var pid1 uint32 = 1
	var pid2 uint32 = 3
	uid1 := m.PM.Alive[pid1].UID
	uid2 := m.PM.Alive[pid2].UID
	uids := []int{uid1, uid2}
	pm := m.GetProcessMeasurements(uids)
	assertEqualsInt(t, "Number of times", 1, len(pm.Times))
	assertEqualsSlice(t, "Values 1", []uint32{pid1 + 1}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{pid2 + 1}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 34
	pMock.Processes[pid2].MemoryUsage = 1024 * 12
	m.measureAndLog(true)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{pid1 + 1, 34}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{pid2 + 1, 12}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 87
	pMock.Processes[pid2].MemoryUsage = 1024 * 21
	m.measureAndLog(false)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 87}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 21}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 44
	pMock.Processes[pid2].MemoryUsage = 1024 * 11
	m.measureAndLog(true)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 87, 44}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 21, 11}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 10
	pMock.Processes[pid2].MemoryUsage = 1024 * 43
	m.measureAndLog(false)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 44, 10}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 11, 43}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 65
	pMock.Processes[pid2].MemoryUsage = 1024 * 56
	m.measureAndLog(true)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 44, 10, 65}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 11, 43, 56}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 87
	pMock.Processes[pid2].MemoryUsage = 1024 * 78
	m.measureAndLog(false)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 44, 65, 87}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 11, 56, 78}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 28
	pMock.Processes[pid2].MemoryUsage = 1024 * 87
	m.measureAndLog(true)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 44, 65, 87, 28}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 11, 56, 78, 87}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 71
	pMock.Processes[pid2].MemoryUsage = 1024 * 17
	m.measureAndLog(false)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{34, 44, 65, 28, 71}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{12, 11, 56, 87, 17}, pm.Memory[uid2])

	time.Sleep(50 * time.Millisecond) // To make time differ
	pMock.Processes[pid1].MemoryUsage = 1024 * 98
	pMock.Processes[pid2].MemoryUsage = 1024 * 89
	m.measureAndLog(true)
	pm = m.GetProcessMeasurements(uids)
	assertEqualsSlice(t, "Values 1", []uint32{44, 65, 28, 71, 98}, pm.Memory[uid1])
	assertEqualsSlice(t, "Values 2", []uint32{11, 56, 87, 17, 89}, pm.Memory[uid2])

	// Get measurement for non existing process
	m.GetProcessMeasurements([]int{12345})
	_, hasElement := pm.Memory[12345]
	assertTrue(t, "No measurement for non-existing process", !hasElement)
}

func TestRemoveOldProcesses(t *testing.T) {
	pMock := proci.GenerateMock(3)
	m := CreateMeasurement(2, 4, 2, 4, pMock)

	m.measureAndLog(true)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 3, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))

	// Kill process with pid 2
	uid := m.PM.Alive[2].UID
	delete(pMock.Processes, 2)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(false)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement := m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(true)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(false)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(true)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(false)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(true)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(false)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 3, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process still available", hasElement)

	time.Sleep(50 * time.Millisecond) // To make time differ
	m.measureAndLog(true)
	m.removeOldProcesses()
	assertEqualsInt(t, "Number of alive processes", 2, len(m.PM.Alive))
	assertEqualsInt(t, "Total number of processes", 2, len(m.PM.All))
	_, hasElement = m.PM.All[uid]
	assertTrue(t, "Removed process has been deleted", !hasElement)
}

func TestMeasurement(t *testing.T) {
	m := CreateMeasurement(2, 4, 200, 3, proci.Proci{})

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
	m := CreateMeasurement(20, 20, 500, 2, proci.Proci{})
	m.Start()

	time.Sleep(3 * time.Second)

	m.Stop()

	t.Log("Size of Fastlogger:", m.FastLogger.NbrRows)
	t.Log("Size of SlowLogger:", m.SlowLogger.NbrRows)
	assertTrue(t, "Size of FastLogger", m.FastLogger.NbrRows > 4 && m.FastLogger.NbrRows < 8)
	assertEqualsInt(t, "Size of SlowLogger", int(m.FastLogger.NbrRows/2), m.SlowLogger.NbrRows)

}

func TestMeasureLoop2(t *testing.T) {
	m := CreateMeasurement(20, 20, 200, 3, proci.Proci{})
	m.Start()

	time.Sleep(3 * time.Second)

	m.Stop()

	t.Log("Size of Fastlogger:", m.FastLogger.NbrRows)
	t.Log("Size of SlowLogger:", m.SlowLogger.NbrRows)
	assertTrue(t, "Size of FastLogger", m.FastLogger.NbrRows > 12 && m.FastLogger.NbrRows < 18)
	assertEqualsInt(t, "Size of SlowLogger", int(m.FastLogger.NbrRows/3), m.SlowLogger.NbrRows)

}
