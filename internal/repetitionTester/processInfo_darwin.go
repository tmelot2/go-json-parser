// +build darwin

package repetitionTester

import (
	"fmt"
	"syscall"
	"unsafe"
)

func GetPageFaultCount() uint64 {
	var info syscall.Rusage

	_, _, errno := syscall.Syscall(syscall.SYS_GETRUSAGE, uintptr(syscall.RUSAGE_SELF), uintptr(unsafe.Pointer(&info)), 0)
	if errno != 0 {
		fmt.Printf("syscall.Syscall error: %v\n", errno)
	}

	return uint64(info.Minflt + info.Majflt)
}
