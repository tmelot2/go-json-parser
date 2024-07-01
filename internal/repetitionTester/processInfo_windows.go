// +build windows

package repetitionTester

import (
	// "fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Copied from MSDN: https://learn.microsoft.com/en-us/windows/win32/api/psapi/ns-psapi-process_memory_counters
type PROCESS_MEMORY_COUNTERS struct {
    cb                         uint32
    PageFaultCount             uint32
    PeakWorkingSetSize         uintptr
    WorkingSetSize             uintptr
    QuotaPeakPagedPoolUsage    uintptr
    QuotaPagedPoolUsage        uintptr
    QuotaPeakNonPagedPoolUsage uintptr
    QuotaNonPagedPoolUsage     uintptr
    PagefileUsage              uintptr
    PeakPagefileUsage          uintptr
}

// Declare syscalls for getting GetProcessMemoryInfo.
var (
	psapi                    = syscall.NewLazyDLL("psapi.dll")
	procGetProcessMemoryInfo = psapi.NewProc("GetProcessMemoryInfo")
)

// Returns value of syscall result for page fault count.
func GetPageFaultCount() uint64 {
	var memCounters PROCESS_MEMORY_COUNTERS
	var err error

	handle := windows.CurrentProcess()
	cb := uint32(unsafe.Sizeof(memCounters))

	r1, _, e1 := syscall.Syscall(
		procGetProcessMemoryInfo.Addr(),
		3,
		uintptr(handle),
		uintptr(unsafe.Pointer(&memCounters)),
		uintptr(cb),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		panic(err)
	}

	return uint64(memCounters.PageFaultCount)
}
