// process unit tests
package main

import (
	"fmt"
	"github.com/midstar/proci"
	"runtime/debug"
	"testing"
)

func assertTrue(t *testing.T, message string, check bool) {
	if !check {
		debug.PrintStack()
		t.Fatal(message)
	}
}

func assertEqualsInt(t *testing.T, message string, expected int, actual int) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %d, Actual: %d", message, expected, actual), expected == actual)
}

func assertEqualsStr(t *testing.T, message string, expected string, actual string) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %s, Actual: %s", message, expected, actual), expected == actual)
}

func TestProcessUpdate(t *testing.T) {
	pMock := proci.GenerateMock(10)
	pMap := NewProcessMap(pMock)
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 10, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 10, len(pMap.Alive))

	p8 := pMap.Alive[8]
	assertEqualsInt(t, "Process 8 PID", 8, int(p8.Pid))
	assertEqualsStr(t, "Process 8 path", "path_8", p8.Path)
	assertEqualsStr(t, "Process 8 command line", "command_line_8", p8.CommandLine)
	assertEqualsInt(t, "Process 8 original MemoryUsage", 1024+1024*8, int(pMock.Processes[8].MemoryUsage)) // Sanity
	assertEqualsInt(t, "Process 8 MaxMemoryEver", 9, int(p8.MaxMemoryEver))
	assertEqualsInt(t, "Process 8 MinMemoryEver", 9, int(p8.MinMemoryEver))
	assertEqualsInt(t, "Process 8 LastMemory", 9, int(p8.LastMemory))

	pMock.Processes[8].MemoryUsage = 1024 * 20
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 10, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 10, len(pMap.Alive))
	assertEqualsInt(t, "Process 8 MaxMemoryEver", 20, int(p8.MaxMemoryEver))
	assertEqualsInt(t, "Process 8 MinMemoryEver", 9, int(p8.MinMemoryEver))
	assertEqualsInt(t, "Process 8 LastMemory", 20, int(p8.LastMemory))

	pMock.Processes[8].MemoryUsage = 1024 * 3
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 10, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 10, len(pMap.Alive))
	assertEqualsInt(t, "Process 8 MaxMemoryEver", 20, int(p8.MaxMemoryEver))
	assertEqualsInt(t, "Process 8 MinMemoryEver", 3, int(p8.MinMemoryEver))
	assertEqualsInt(t, "Process 8 LastMemory", 3, int(p8.LastMemory))

	pMock.Processes[8].MemoryUsage = 1024 * 4
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 10, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 10, len(pMap.Alive))
	assertEqualsInt(t, "Process 8 MaxMemoryEver", 20, int(p8.MaxMemoryEver))
	assertEqualsInt(t, "Process 8 MinMemoryEver", 3, int(p8.MinMemoryEver))
	assertEqualsInt(t, "Process 8 LastMemory", 4, int(p8.LastMemory))

	delete(pMock.Processes, 8)
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 10, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 9, len(pMap.Alive))
	_, hasP8 := pMap.Alive[8]
	assertTrue(t, "PID P8 dead", !hasP8)

	p3 := pMap.Alive[3]
	assertEqualsStr(t, "Process 3 path", "path_3", p3.Path)
	pMock.Processes[3].Path = "new_path_3"
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 11, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 9, len(pMap.Alive))
	assertTrue(t, "New PID 3 differs from old", p3.UID != pMap.Alive[3].UID)

	p2 := pMap.Alive[2]
	pMock.Processes[2].Path = ""
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 11, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 8, len(pMap.Alive))
	assertTrue(t, "New PID 2 is dead", !p2.IsAlive)

	pMock.Processes[34] = &proci.ProcessMock{
		Pid:               34,
		Path:              fmt.Sprintf("path_%d", 34),
		CommandLine:       fmt.Sprintf("command_line_%d", 34),
		MemoryUsage:       uint64(1024 + 34*1024),
		DoFailPath:        false,
		DoFailCommandLine: false,
		DoFailMemoryUsage: false}

	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 12, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 9, len(pMap.Alive))
	assertEqualsStr(t, "Process 34 path", "path_34", pMap.Alive[34].Path)

	pMock.Processes[22] = &proci.ProcessMock{
		Pid:               22,
		Path:              "",
		CommandLine:       fmt.Sprintf("command_line_%d", 22),
		MemoryUsage:       uint64(1024 + 22*1024),
		DoFailPath:        false,
		DoFailCommandLine: false,
		DoFailMemoryUsage: false}
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 12, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 9, len(pMap.Alive))
	_, hasPid22 := pMap.Alive[22]
	assertTrue(t, "Process 22 shall be ignored", !hasPid22)

	pMock.Processes[23] = &proci.ProcessMock{
		Pid:               23,
		Path:              fmt.Sprintf("path_%d", 23),
		CommandLine:       fmt.Sprintf("command_line_%d", 23),
		MemoryUsage:       uint64(1024 + 23*1024),
		DoFailPath:        true,
		DoFailCommandLine: false,
		DoFailMemoryUsage: false}
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 12, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 9, len(pMap.Alive))
	_, hasPid23 := pMap.Alive[23]
	assertTrue(t, "Process 23 shall be ignored", !hasPid23)

	pMock.Processes[24] = &proci.ProcessMock{
		Pid:               24,
		Path:              fmt.Sprintf("path_%d", 24),
		CommandLine:       fmt.Sprintf("command_line_%d", 24),
		MemoryUsage:       uint64(1024 + 24*1024),
		DoFailPath:        false,
		DoFailCommandLine: true,
		DoFailMemoryUsage: true}
	pMap.Update()
	assertEqualsInt(t, "Length of pMap.All", 13, len(pMap.All))
	assertEqualsInt(t, "Length of pMap.Alive", 10, len(pMap.Alive))
	assertEqualsStr(t, "Process 24 has empty CommandLine", "", pMap.Alive[24].CommandLine)
	assertEqualsInt(t, "Process 24 has no Memory", 0, int(pMap.Alive[24].LastMemory))

	////////////////////////////////////////////////////////////////////////////
	// Physical memory
	assertEqualsInt(t, "Size of pMap.TotalPhys", 4*1024*1024, int(pMap.TotalPhys))
	assertEqualsInt(t, "Size of pMap.MaxPhysEver", 2*1024*1024, int(pMap.MaxPhysEver))
	assertEqualsInt(t, "Size of pMap.MinPhysEver", 2*1024*1024, int(pMap.MinPhysEver))
	assertEqualsInt(t, "Size of pMap.LastPhys", 2*1024*1024, int(pMap.LastPhys))

	pMock.MemStatus.AvailPhys = 3 * 1024 * 1024 * 1024 // Bytes
	pMap.Update()
	assertEqualsInt(t, "Size of pMap.TotalPhys", 4*1024*1024, int(pMap.TotalPhys))
	assertEqualsInt(t, "Size of pMap.MaxPhysEver", 2*1024*1024, int(pMap.MaxPhysEver))
	assertEqualsInt(t, "Size of pMap.MinPhysEver", 1*1024*1024, int(pMap.MinPhysEver))
	assertEqualsInt(t, "Size of pMap.LastPhys", 1*1024*1024, int(pMap.LastPhys))

	pMock.MemStatus.AvailPhys = 1 * 1024 * 1024 * 1024 // Bytes
	pMap.Update()
	assertEqualsInt(t, "Size of pMap.TotalPhys", 4*1024*1024, int(pMap.TotalPhys))
	assertEqualsInt(t, "Size of pMap.MaxPhysEver", 3*1024*1024, int(pMap.MaxPhysEver))
	assertEqualsInt(t, "Size of pMap.MinPhysEver", 1*1024*1024, int(pMap.MinPhysEver))
	assertEqualsInt(t, "Size of pMap.LastPhys", 3*1024*1024, int(pMap.LastPhys))

	pMock.DoFailMemStatus = true
	pMap.Update()
	assertEqualsInt(t, "Size of pMap.TotalPhys", 4*1024*1024, int(pMap.TotalPhys))
	assertEqualsInt(t, "Size of pMap.MaxPhysEver", 3*1024*1024, int(pMap.MaxPhysEver))
	assertEqualsInt(t, "Size of pMap.MinPhysEver", 1*1024*1024, int(pMap.MinPhysEver))
	assertEqualsInt(t, "Size of pMap.LastPhys", 0*1024*1024, int(pMap.LastPhys))

}
