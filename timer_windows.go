// +build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Declare asm functions that return CPU timer values.
func ReadCpuTimerStart() int64
func ReadCpuTimerEnd()   int64

// Declare syscalls for getting QueryPerformanceCounter
var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	procQueryPerformanceFrequency = kernel32.NewProc("QueryPerformanceFrequency")
	procQueryPerformanceCounter   = kernel32.NewProc("QueryPerformanceCounter")
)

// Returns result of syscall QueryPerformanceFrequency()
func GetOSTimerFreq() (int64, error) {
	var freq int64
	var err error
	r1, _, e1 := syscall.Syscall(
		procQueryPerformanceFrequency.Addr(),
		1,
		uintptr(unsafe.Pointer(&freq)),
		0,
		0,
	)

	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return freq, err
	}
	return freq, err
}

// Returns result of syscall QueryPerformanceCounter()
func ReadOSTimer() (int64, error) {
	var osTimer int64
	var err error
	r1, _, e1 := syscall.Syscall(
		procQueryPerformanceCounter.Addr(),
		1,
		uintptr(unsafe.Pointer(&osTimer)),
		0,
		0,
	)

	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
		return osTimer, err
	}
	return osTimer, err
}

// Prints read, measurement, & guess of CPU frequency & related data.
func PrintTimerStats() {
	fmt.Println("[CPU timer stats]")

	// Setup
	millisecondsToWait := int64(10)
	width := 16 // Output width
	p := message.NewPrinter(language.English) // For printing large numbers with commas

	// Get OS timer frequency
	osFreq, err := GetOSTimerFreq()
	if err != nil {
		fmt.Println("Error getting OS timer frequency:", err)
		return
	}
	// In nanoseconds per second
	p.Printf("OS Timer Frequency [reported]:          %*d\n", width, osFreq)

	cpuStart   := ReadCpuTimerStart()
	osStart, _ := ReadOSTimer()

	var osEnd int64
	var osElapsed int64
	osWaitTime := osFreq * millisecondsToWait / 1000
	for osElapsed < osWaitTime {
		osEnd, _ = ReadOSTimer()
		osElapsed = osEnd - osStart
	}
	cpuEnd := ReadCpuTimerEnd()
	cpuElapsed := cpuEnd - cpuStart

	cpuFreq := int64(0)
	if osElapsed > 0 {
		cpuFreq = osFreq * cpuElapsed / osElapsed
	}

	// p.Printf(  "OS Timer:      %*d -> %*d = %*d elapsed\n", width, osStart, width, osEnd, width, osElapsed)
	p.Printf("OS Timer:                               %*d elapsed\n", width, osElapsed)
	p.Printf("OS Seconds (elapsed/freq):                   %*.4f\n", width, float64(osElapsed) / float64(osFreq))

	// p.Printf(  "CPU timer:     %*d -> %*d = %*d elapsed\n", width, cpuStart, width, cpuEnd, width, cpuElapsed)
	p.Printf("CPU timer:                              %*d elapsed\n", width, cpuElapsed)
	p.Printf("CPU freq (guessed):                     %*d\n", width, cpuFreq)

}
