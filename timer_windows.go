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
func ReadCPUTimer()      uint64

// Declare syscalls for getting QueryPerformanceCounter
var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	procQueryPerformanceFrequency = kernel32.NewProc("QueryPerformanceFrequency")
	procQueryPerformanceCounter   = kernel32.NewProc("QueryPerformanceCounter")
)

// Returns result of syscall QueryPerformanceFrequency()
// Returns 10,000,000 on my amd64 Windows machine
func GetOSTimerFreq() (uint64, error) {
	var freq uint64
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
func ReadOSTimer() (uint64, error) {
	var osTimer uint64
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
func EstimateCPUTimerFreq(printDebug bool) uint64 {
	// Setup
	millisecondsToWait := uint64(100)
	width := 20 // Output width
	p := message.NewPrinter(language.English) // For printing large numbers with commas

	// Get OS timer frequency
	osFreq, err := GetOSTimerFreq()
	if err != nil {
		fmt.Println("Error getting OS timer frequency:", err)
		return 0
	}
	// In nanoseconds per second
	if printDebug {
		p.Printf("OS Timer Frequency [reported]: %*d\n", width, osFreq)
	}

	cpuStart   := ReadCPUTimer()
	osStart, _ := ReadOSTimer()

	var osEnd uint64
	var osElapsed uint64
	osWaitTime := osFreq * millisecondsToWait / 1000
	for osElapsed < osWaitTime {
		osEnd, _ = ReadOSTimer()
		osElapsed = osEnd - osStart
	}
	cpuEnd := ReadCPUTimer()
	cpuElapsed := cpuEnd - cpuStart

	cpuFreq := uint64(0)
	if osElapsed > 0 {
		cpuFreq = osFreq * cpuElapsed / osElapsed
	}

	if printDebug {
		p.Printf("OS Timer:                      %*d elapsed\n", width, osElapsed)
		p.Printf("OS Seconds (elapsed/freq):          %*.4f\n", width, float64(osElapsed) / float64(osFreq))

		p.Printf("CPU timer:                     %*d elapsed\n", width, cpuElapsed)
		p.Printf("CPU freq (guessed):            %*d\n", width, cpuFreq)
	}

	return uint64(cpuFreq)
}
